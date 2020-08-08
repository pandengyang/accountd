package controllers

import (
	"accountd/datamodels"
	"accountd/services"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pandengyang/utils/CollectionJSON"
	"github.com/pandengyang/utils/StringUtils"
	"io/ioutil"
	"time"
)

var (
	pTokenTemplateStr *string
)

type TokenController struct {
	Service        services.TokenService
	AccountService services.AccountService

	PrivateKeyPathname string
	PrivateKey         *ecdsa.PrivateKey
}

func init() {
	var err error
	var content []byte

	if content, err = ioutil.ReadFile("web/views/CollectionJSON/v1/token.min.json"); err != nil {
		panic("read template error!")
	}

	contentStr := string(content)
	pTokenTemplateStr = &contentStr
}

func (c *TokenController) BeforeActivation(ba mvc.BeforeActivation) {
	keyData, err := ioutil.ReadFile(c.PrivateKeyPathname)
	if err != nil {
		panic(fmt.Errorf("read file %s error: %v", c.PrivateKeyPathname, err))
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(keyData)
	if err != nil {
		panic(fmt.Errorf("parse %s error: %v", c.PrivateKeyPathname, err))
	}
	c.PrivateKey = privateKey

	ba.Handle("POST", "/", "Post")
	ba.Handle("Delete", "/", "Delete")
}

func (c *TokenController) Post(ctx iris.Context) mvc.Result {
	var err error

	var datas CollectionJSON.Datas
	var account datamodels.Account
	var password string
	var refreshTokenExists bool

	var accessToken string
	var refreshToken string
	var tokenJson string

	err = ctx.ReadJSON(&datas)
	if err != nil {
		err = fmt.Errorf("ReadJSON err: %v", err)

		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	for _, value := range datas.Data {
		switch value.Name {
		case "nickname":
			account.Nickname = value.Value

		case "password":
			password = value.Value

		case "phone":
			account.Phone = value.Value

		case "verification_code":
			account.VerificationCode = value.Value

		case "refresh_token":
			refreshToken = value.Value
		}
	}

	if account.Phone != "" {
		if account, err = c.AccountService.SelectAuthByPhone(account.Phone); err != nil {
			if err == sql.ErrNoRows {
				return mvc.Response{
					Code:        iris.StatusUnauthorized,
					ContentType: "text/plain",
					Text:        fmt.Sprintf("%v", err),
				}
			}

			return mvc.Response{
				Code:        iris.StatusInternalServerError,
				ContentType: "text/plain",
				Text:        fmt.Sprintf("%v", err),
			}
		}
	} else if account.Nickname != "" {
		account, err = c.AccountService.SelectAuthByNickname(account.Nickname)
		if err != nil {
			if err == sql.ErrNoRows {
				return mvc.Response{
					Code:        iris.StatusUnauthorized,
					ContentType: "text/plain",
					Text:        fmt.Sprintf("%v", errors.New("invalid account or password")),
				}
			}

			return mvc.Response{
				Code:        iris.StatusInternalServerError,
				ContentType: "text/plain",
				Text:        fmt.Sprintf("%v", err),
			}
		}

		saltedHashedPassword := StringUtils.Sha256PasswdSalt(password, account.Salt)
		if account.Password != saltedHashedPassword {
			return mvc.Response{
				Code:        iris.StatusUnauthorized,
				ContentType: "text/plain",
				Text:        fmt.Sprintf("%v", errors.New("invalid account or password")),
			}
		}
	} else if refreshToken != "" {
		if refreshTokenExists, err = c.Service.RefreshTokenExists(refreshToken); err != nil {
			return mvc.Response{
				Code:        iris.StatusInternalServerError,
				ContentType: "text/plain",
				Text:        fmt.Sprintf("%v", err),
			}
		}

		if !refreshTokenExists {
			return mvc.Response{
				Code:        iris.StatusUnauthorized,
				ContentType: "text/plain",
				Text:        fmt.Sprintf("%v", errors.New("refresh token expires")),
			}
		}
	} else {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("account info is empty")),
		}
	}

	now := time.Now().Unix()
	jti := fmt.Sprintf("%d-%d-%s", account.Id, now, StringUtils.GetRandomString(datamodels.JTI_RANDOM_LEN))
	claims := &jwt.MapClaims{
		"aud":      datamodels.AUDIENCE,
		"exp":      now + int64(datamodels.EFFECTIVE_TIME),
		"jti":      jti,
		"iat":      now,
		"iss":      datamodels.ISSUER,
		"nbf":      now,
		"sub":      fmt.Sprintf("%d", account.Id),
		"nickname": account.Nickname,
	}

	if accessToken, err = jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(c.PrivateKey); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if refreshToken == "" {
		refreshToken = StringUtils.Md5(jti)
		if _, err = c.Service.InsertRefreshToken(refreshToken); err != nil {
			return mvc.Response{
				Code:        iris.StatusInternalServerError,
				ContentType: "text/plain",
				Text:        fmt.Sprintf("%v", err),
			}
		}
	}

	token := datamodels.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if tokenJson, err = CollectionJSON.Item(token, pTokenTemplateStr); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code: iris.StatusCreated,
		Text: tokenJson,
	}
}

func (c *TokenController) Delete(ctx iris.Context) mvc.Result {
	var err error

	ctx.Values().Get("user_id")
	jti = ctx.Values.Get("jti")

	/* 删除 refresh token */
	if _, err = c.Service.DeleteRefreshToken(refreshToken); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	/* 吊销 access token */
	if _, err = c.Service.RevokeAccessToken(jti); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code: iris.StatusOK,
		Text: "ok",
	}
}
