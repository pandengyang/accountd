package main

import (
	"github.com/gavv/httpexpect"
	"github.com/kataras/iris/v12/httptest"
	"github.com/pandengyang/utils/CollectionJSON"
	"testing"
)

func TestNewApp(t *testing.T) {
	var response *httpexpect.Response

	app := newApp()
	e := httptest.New(t, app)

	/*
		datas := CollectionJSON.Datas{
			Data: []CollectionJSON.Data{
				CollectionJSON.Data{"phone", "18612466738", ""},
			},
		}
		response = e.POST("/api/verificationcodes").WithJSON(datas).Expect().Status(httptest.StatusCreated)
		t.Log(response.Body())
	*/

	/*
		datas := CollectionJSON.Datas{
			Data: []CollectionJSON.Data{
				CollectionJSON.Data{"nickname", "PanDengyang", ""},
				CollectionJSON.Data{"phone", "18612466738", ""},
				CollectionJSON.Data{"verification_code", "5438", ""},
				CollectionJSON.Data{"password", "123456", ""},
			},
		}
		response = e.POST("/api/accounts").WithJSON(datas).Expect().Status(httptest.StatusCreated)
		t.Log(response.Body())
	*/

	datas := CollectionJSON.Datas{
		Data: []CollectionJSON.Data{
			CollectionJSON.Data{"nickname", "PanDengyang", ""},
			CollectionJSON.Data{"password", "123456", ""},
		},
	}
	response = e.POST("/api/tokens").WithJSON(datas).Expect().Status(httptest.StatusCreated)
	t.Log(response.Body())
}
