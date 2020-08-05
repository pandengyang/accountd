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
	"time"
)

type VerificationCodeController struct {
	Service     services.VerificationCodeService
	Middlewares []iris.Handler
}

func (c *VerificationCodeController) BeforeActivation(ba mvc.BeforeActivation) {
	ba.Handle("POST", "/", "Post")
}

func (c *VerificationCodeController) Post(ctx iris.Context) mvc.Result {
	var err error

	var datas CollectionJSON.Datas
	var vc datamodels.VerificationCode
	var vcOld datamodels.VerificationCode

	var insertedPhone string

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
			vc.Phone = value.Value
		}
	}

	if vc.Phone == "" {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("verification code info is empty")),
		}
	}

	if vcOld, err = c.Service.SelectByPhone(vc.Phone); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	if time.Now().Unix() <= vcOld.SentAt+datamodels.VERIFICATION_CODE_SEND_INTERVAL {
		return mvc.Response{
			Code:        iris.StatusBadRequest,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", errors.New("sending verification codes too frequently")),
		}
	}

	vc.VerificationCode = StringUtils.GetVerificationCode(datamodels.VERIFICATION_CODE_LEN)
	vc.SentAt = time.Now().Unix()

	// 发送验证码

	if insertedPhone, err = c.Service.Insert(&vc); err != nil {
		return mvc.Response{
			Code:        iris.StatusInternalServerError,
			ContentType: "text/plain",
			Text:        fmt.Sprintf("%v", err),
		}
	}

	return mvc.Response{
		Code: iris.StatusCreated,
		Text: fmt.Sprintf("/verificationcodes/%s", insertedPhone),
	}
}
