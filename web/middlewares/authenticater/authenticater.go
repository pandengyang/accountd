package authenticater

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	irisjwt "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
)

const (
	CONTEXT_KEY = "token"
)

var (
	publicKey *ecdsa.PublicKey
)

func New(publicKeyPathname string) *irisjwt.Middleware {
	var err error

	keyData, err := ioutil.ReadFile(publicKeyPathname)
	if err != nil {
		panic(fmt.Errorf("read file %s error: %v", publicKeyPathname, err))
	}

	publicKey, err := jwt.ParseECPublicKeyFromPEM(keyData)
	if err != nil {
		panic(fmt.Errorf("parse %s error: %v", publicKeyPathname, err))
	}

	return irisjwt.New(irisjwt.Config{
		ValidationKeyGetter: func(token *irisjwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		ContextKey:    CONTEXT_KEY,
		SigningMethod: irisjwt.SigningMethodES256,
	})
}

func ExtractClaims(ctx iris.Context) {
	token := ctx.Values().Get(CONTEXT_KEY).(*irisjwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	ctx.Values().Set("user_id", claims["sub"])
	ctx.Values().Set("jwt_id", claims["jti"])

	ctx.Next()
}

/* 检测 jwt 是否被吊销 */
func CheckTokenRevoked(ctx iris.Context) {
	revoked := false
	if revoked {
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString("token revoked")

		return
	}

	ctx.Next()
}
