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
	irisjwt "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"os"
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

	Db             *sql.DB
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

	/* Database */
	if err = initDatabase(app); err != nil {
		panic(err)
	}

	/* Redis 连接池 */

	/* 配置 mvc */
	mvc.Configure(app.Party("/api/accounts"), accounts)
	mvc.Configure(app.Party("/api/auth"), auth)

	return app
}

func initDatabase(app *iris.Application) (err error) {
	dbConfig := app.ConfigurationReadOnly().GetOther()["Database"].(map[string]interface{})
	dbDriver := dbConfig["driver"].(string)
	dbHost := dbConfig["host"].(string)
	dbPort := dbConfig["port"].(string)
	dbUser := dbConfig["user"].(string)
	dbPassword := dbConfig["password"].(string)
	dbDatabase := dbConfig["database"].(string)

	/* SQL 连接池 */
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	Db, err = sql.Open(dbDriver, dataSourceName)

	return err
}

func accounts(app *mvc.Application) {
	middlewares := []iris.Handler{
		authMiddleware.Serve,
		authenticater.ExtractClaims,
		authenticater.CheckTokenRevoked,
	}

	/* 数据仓库 */
	persistenceRepo := repositories.NewAccountMySQLRepository(Db)

	/* 服务 */
	userService := services.NewAccountService(persistenceRepo, nil)
	app.Register(userService)

	app.Handle(&controllers.AccountController{
		Middlewares: middlewares,
	})
}

func auth(app *mvc.Application) {
	middlewares := []iris.Handler{
		authMiddleware.Serve,
		authenticater.ExtractClaims,
		authenticater.CheckTokenRevoked,
	}

	/* 数据仓库 */
	persistenceRepo := repositories.NewAccountMySQLRepository(Db)

	/* 服务 */
	userService := services.NewAccountService(persistenceRepo, nil)
	app.Register(userService)

	app.Handle(&controllers.AuthController{
		PrivateKeyPathname: PRIKEY,
		Middlewares:        middlewares,
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
