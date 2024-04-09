package integrationtest_test

import (
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

const (
	host     = "192.168.56.1:8080"
	basePath = "http://" + host
)

func TestGetBanner(t *testing.T) {
	var userToken string
	var admToken string

	MustDo(
		Get(basePath+"/get_token/1"),
		Store().Response().Body().JSON().JQ(".token").In(&admToken),
	)

	MustDo(
		Get(basePath+"/get_token/0"),
		Store().Response().Body().JSON().JQ(".token").In(&userToken),
	)

	Test(t,
		Description("GetBanner Success"),
		Send().Headers("Authorization").Add("Bearer "+userToken),
		Get(basePath+"/user_banner?tag_id=1&feature_id=1"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("title"),
		Expect().Body().JSON().Contains("text"),
		Expect().Body().JSON().Contains("url"),
	)

	Test(t,
		Description("GetBanner fail, banner is inactive"),
		Send().Headers("Authorization").Add("Bearer "+userToken),
		Get(basePath+"/user_banner?tag_id=3&feature_id=1"),
		Expect().Status().Equal(http.StatusNotFound),
	)

	Test(t,
		Description("GetBanner admin success, banner is inactive"),
		Send().Headers("Authorization").Add("Bearer "+admToken),
		Get(basePath+"/user_banner?tag_id=3&feature_id=1"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Contains("title"),
		Expect().Body().JSON().Contains("text"),
		Expect().Body().JSON().Contains("url"),
	)

	Test(t,
		Description("GetBanner admin fail, incorrect data"),
		Send().Headers("Authorization").Add("Bearer "+admToken),
		Get(basePath+"/user_banner?tag_id=fsdf&feature_id=1"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().Contains("error"),
	)

	Test(t,
		Description("GetBanner fail, unauthorized"),
		Get(basePath+"/user_banner?tag_id=fsdf&feature_id=1"),
		Expect().Status().Equal(http.StatusUnauthorized),
	)
}
