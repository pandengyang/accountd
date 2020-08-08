package main

import (
	"accountd/repositories"
	"accountd/services"
	"accountd/web/controllers"
	"accountd/web/middlewares/authenticater"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	irisjwt "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"os"
	"time"
)

const (
	PUBKEY = "env/working/accountd_pubkey.pem"
	PRIKEY = "env/working/accountd_prikey.pem"

	LOGFILE = "env/working/log/accountd.log"
	PIDFILE = "env/working/accountd.pid"
)

var (
	/* flags */
	ip   = flag.String("ip", "", "The ip address where the server listens to")
	port = flag.String("port", "", "The port which the server listens to")

	SqlDb   *sql.DB
	CacheDb *redis.Pool

	authMiddleware *irisjwt.Middleware
)

func main() {
	/* 解析标志 */
	if err := initOption(); err != nil {
		panic(err)
	}

	app := newApp()

	/* 处理信号与日志 */
	initLog(app)
	wrPid()
	defer rmPid()

	app.Run(iris.Addr(fmt.Sprintf("%s:%s", *ip, *port)))
}

func initOption() (err error) {
	flag.Parse()

	if "" == *ip || "" == *port {
		err = fmt.Errorf("ip or port is empty!")

		return err
	}

	return nil
}

func newApp() (app *iris.Application) {
	var err error

	authMiddleware = authenticater.New(PUBKEY)

	app = iris.New()

	app.Use(recover.New())
	app.Use(logger.New())

	/* 配置 */
	app.Configure(iris.WithConfiguration(iris.TOML("env/working/config.tml")))

	/* SQL 数据库 */
	if err = initSqlDb(app); err != nil {
		panic(err)
	}

	/* 缓存 */
	if err = initCacheDb(app); err != nil {
		panic(err)
	}

	/* 配置 mvc */
	mvc.Configure(app.Party("/api/verificationcodes"), verificationCodes)
	mvc.Configure(app.Party("/api/accounts"), accounts)
	mvc.Configure(app.Party("/api/tokens"), tokens)

	return app
}

func initSqlDb(app *iris.Application) (err error) {
	dbConfig := app.ConfigurationReadOnly().GetOther()["SqlDb"].(map[string]interface{})
	dbDriver := dbConfig["driver"].(string)
	dbHost := dbConfig["host"].(string)
	dbPort := dbConfig["port"].(string)
	dbUser := dbConfig["user"].(string)
	dbPassword := dbConfig["password"].(string)
	dbDatabase := dbConfig["database"].(string)

	/* SQL 连接池 */
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	SqlDb, err = sql.Open(dbDriver, dataSourceName)

	return err
}

func initCacheDb(app *iris.Application) (err error) {
	dbConfig := app.ConfigurationReadOnly().GetOther()["CacheDb"].(map[string]interface{})
	dbHost := dbConfig["host"].(string)
	dbPort := dbConfig["port"].(string)
	dbPassword := dbConfig["password"].(string)
	dbDatabase := dbConfig["database"].(int64)

	CacheDb = &redis.Pool{
		MaxActive:   1000,              // 最激活闲连接数
		MaxIdle:     100,               // 最大空闲连接数
		IdleTimeout: 300 * time.Second, // 空闲连接关闭超时
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", dbHost, dbPort),
				redis.DialPassword(dbPassword),
				redis.DialDatabase(int(dbDatabase)),
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(5*time.Second),
				redis.DialWriteTimeout(5*time.Second))
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}

	return err
}

func verificationCodes(app *mvc.Application) {
	/* 数据仓库 */
	cacheRepo := repositories.NewVerificationCodeRedisRepository(CacheDb)

	/* 服务 */
	vcService := services.NewVerificationCodeService(nil, cacheRepo)
	app.Register(vcService)

	app.Handle(&controllers.VerificationCodeController{})
}

func accounts(app *mvc.Application) {
	middlewares := []iris.Handler{
		authMiddleware.Serve,
		authenticater.ExtractClaims,
		authenticater.CheckTokenRevoked,
	}

	/* 数据仓库 */
	persistenceRepo := repositories.NewAccountMySQLRepository(SqlDb)
	cacheRepo := repositories.NewVerificationCodeRedisRepository(CacheDb)

	/* 服务 */
	accountService := services.NewAccountService(persistenceRepo, nil)
	app.Register(accountService)

	vcService := services.NewVerificationCodeService(nil, cacheRepo)
	app.Register(vcService)

	app.Handle(&controllers.AccountController{
		Middlewares: middlewares,
	})
}

func tokens(app *mvc.Application) {
	/* 数据仓库 */
	persistenceRepo := repositories.NewAccountMySQLRepository(SqlDb)
	cacheRepo := repositories.NewTokenRedisRepository(CacheDb)

	/* 服务 */
	tokenService := services.NewTokenService(nil, cacheRepo)
	app.Register(tokenService)

	accountService := services.NewAccountService(persistenceRepo, nil)
	app.Register(accountService)

	app.Handle(&controllers.TokenController{
		PrivateKeyPathname: PRIKEY,
	})
}

func wrPid() (err error) {
	var pidFile *os.File

	if pidFile, err = os.OpenFile(PIDFILE, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
		err = fmt.Errorf("failed to open pid file: %v", err)

		return err
	}
	defer pidFile.Close()

	pidFile.WriteString(fmt.Sprintf("%d", os.Getpid()))
	if err != nil {
		err = fmt.Errorf("failed to write pid file: %v", err)

		return err
	}

	return nil
}

func rmPid() (err error) {
	if err = os.Remove(PIDFILE); err != nil {
		err = fmt.Errorf("failed to open pid file: %v", err)

		return err
	}

	return nil
}

func initLog(app *iris.Application) (err error) {
	var loggerFile *os.File

	if loggerFile, err = os.OpenFile(LOGFILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		err = fmt.Errorf("failed to open log file: %v", err)

		return err
	}

	app.Logger().SetOutput(loggerFile)

	return nil
}
