package controllers

import (
	"accountd/datamodels"
	"accountd/services"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pandengyang/utils/CollectionJSON"
	"github.com/pandengyang/utils/StringUtils"
	"io/ioutil"
	"strconv"
)

var (
	pAccountTemplateStr *string
)

type AccountController struct {
	Service     services.AccountService
	VcService   services.VerificationCodeService
	Middlewares []iris.Handler
}

func init() {
	var err error
	var content []byte

	if content, err = ioutil.ReadFile("web/views/CollectionJSON/v1/account.min.json"); err != nil {
		panic("read template error!")
	}

	contentStr := string(content)
	pAccountTemplateStr = &contentStr
}

func (c *AccountController) BeforeActivation(ba mvc.BeforeActivation) {
	ba.Handle("POST", "/", "Post")
	ba.Handle("DELETE", "/{id:int64}", "Delete")

	ba.Handle("PUT", "/{id:int64}/nickname", "PutNickname")
	ba.Handle("PUT", "/{id:int64}/phone", "PutPhone")
	ba.Handle("PUT", "/{id:int64}/password", "PutPassword")

	ba.Handle("GET", "/", "GetAllPerPage")
	ba.Handle("GET", "/all", "GetAll")
}

func (c *AccountController) Post(ctx iris.Context) mvc.Result {
	var err error

	var datas CollectionJSON.Datas
	var account datamodels.Account
	var vc datamodels.VerificationCode

	var insertedId int64

	if err = ctx.ReadJSON(&datas); err != nil {
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

		case "phone":
			account.Phone = value.Value

		case "verification_code":
			account.VerificationCode = value.Value

		case "password":
			account.Password = value.Value
		}
	}

	if account.Nickname == "" || account.Phone == "" || account.VerificationCode == "" || account.Password == "" {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("account info is empty")),
		}
	}

	if vc, err = c.VcService.SelectByPhone(account.Phone); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if vc.VerificationCode != account.VerificationCode {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("verification code error")),
		}
	}

	account.Salt = StringUtils.GetRandomString(datamodels.SALT_LEN)
	account.Password = StringUtils.Sha256PasswdSalt(account.Password, account.Salt)

	if insertedId, err = c.Service.Insert(&account); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code: iris.StatusCreated,
		Text: fmt.Sprintf("/accounts/%d", insertedId),
	}
}

func (c *AccountController) Delete(ctx iris.Context) mvc.Result {
	var err error

	var id int64
	var rowsAffected int64

	if id, err = ctx.Params().GetInt64("id"); err != nil {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if rowsAffected, err = c.Service.Delete(id); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code:        iris.StatusOK,
		ContentType: "text/plain",
		Text:        fmt.Sprintf("%d", rowsAffected),
	}
}

func (c *AccountController) PutNickname(ctx iris.Context) mvc.Result {
	var err error

	var id int64
	var datas CollectionJSON.Datas
	var account datamodels.Account

	var rowsAffected int64

	if id, err = ctx.Params().GetInt64("id"); err != nil {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if err = ctx.ReadJSON(&datas); err != nil {
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
		}
	}

	if account.Nickname == "" {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("account info is empty")),
		}
	}

	if rowsAffected, err = c.Service.UpdateNickname(id, &account); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code:        iris.StatusOK,
		ContentType: "text/plain",
		Text:        fmt.Sprintf("%d", rowsAffected),
	}
}

func (c *AccountController) PutPhone(ctx iris.Context) mvc.Result {
	var err error

	var id int64
	var datas CollectionJSON.Datas
	var account datamodels.Account

	var rowsAffected int64

	if id, err = ctx.Params().GetInt64("id"); err != nil {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if err = ctx.ReadJSON(&datas); err != nil {
		err = fmt.Errorf("ReadJSON err: %v", err)

		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	for _, value := range datas.Data {
		switch value.Name {
		case "phone":
			account.Phone = value.Value

		case "verification_code":
			account.VerificationCode = value.Value
		}
	}

	if account.Phone == "" || account.VerificationCode == "" {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("account info is empty")),
		}
	}

	if rowsAffected, err = c.Service.UpdatePhone(id, &account); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code:        iris.StatusOK,
		ContentType: "text/plain",
		Text:        fmt.Sprintf("%d", rowsAffected),
	}
}

func (c *AccountController) PutPassword(ctx iris.Context) mvc.Result {
	var err error

	var id int64
	var datas CollectionJSON.Datas
	var account datamodels.Account

	var rowsAffected int64

	if id, err = ctx.Params().GetInt64("id"); err != nil {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if err = ctx.ReadJSON(&datas); err != nil {
		err = fmt.Errorf("ReadJSON err: %v", err)

		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	for _, value := range datas.Data {
		switch value.Name {
		case "phone":
			account.Phone = value.Value

		case "verification_code":
			account.VerificationCode = value.Value

		case "password":
			account.Password = value.Value
		}
	}

	if account.Password == "" || account.Phone == "" || account.VerificationCode == "" {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("account info is empty")),
		}
	}

	account.Salt = StringUtils.GetRandomString(datamodels.SALT_LEN)
	account.Password = StringUtils.Sha256PasswdSalt(account.Password, account.Salt)

	if rowsAffected, err = c.Service.UpdatePassword(id, &account); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code:        iris.StatusOK,
		ContentType: "text/plain",
		Text:        fmt.Sprintf("%d", rowsAffected),
	}
}

func (c *AccountController) GetAll(ctx iris.Context) mvc.Result {
	var err error

	var accountsJson string
	var accounts []interface{}
	var total int64

	if accounts, total, err = c.Service.SelectAll(); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if accountsJson, err = CollectionJSON.Items(accounts, total, pAccountTemplateStr); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code:        iris.StatusOK,
		ContentType: "text/plain",
		Text:        accountsJson,
	}
}

func (c *AccountController) GetAllPerPage(ctx iris.Context) mvc.Result {
	var err error

	var page int64
	var pageSize int64

	var accountsJson string
	var accounts []interface{}
	var total int64

	pageStr := ctx.URLParam("p")
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil {
			return mvc.Response{
				Code:        iris.StatusBadRequest,
				ContentType: "text/plain",
				Text:        "page error",
			}
		}

		pageSizeStr := ctx.URLParam("ps")
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil {
			return mvc.Response{
				Code:        iris.StatusBadRequest,
				ContentType: "text/plain",
				Text:        "page size error",
			}
		}
	}

	if accounts, total, err = c.Service.SelectAllPerPage(page, pageSize); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if accountsJson, err = CollectionJSON.Items(accounts, total, pAccountTemplateStr); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code:        iris.StatusOK,
		ContentType: "text/plain",
		Text:        accountsJson,
	}
}
