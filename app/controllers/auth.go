package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/url"
)

type (
	Auth struct {
	}
	LoginParam struct {
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required_without=UserApiKey"`
		UserApiKey string `json:"user_api_key" binding:"omitempty,min=1"`
	}
	UserInfo struct {
		Email      string
		UserKey    string
		UserApiKey string
	}
)
type CloudflareResponse struct {
	Msg     interface{} `json:"msg"`
	Request struct {
		Act string `json:"act"`
	} `json:"request"`
	Response struct {
		CloudflareEmail string      `json:"cloudflare_email"`
		UniqueId        interface{} `json:"unique_id"`
		UserApiKey      string      `json:"user_api_key"`
		UserKey         string      `json:"user_key"`
	} `json:"response"`
	Result string `json:"result"`
}

func (i *Auth) Login(c *gin.Context) {
	//https://api.cloudflare.com/host-gw.html
	params := &LoginParam{}
	err := c.ShouldBind(&params)
	if err != nil {
		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	userInfo := &util.User{}
	if params.UserApiKey == "" {
		if pkg.Conf.Cloudflare.HostKey == "" {
			util.JSON(c, "error", "HostKey有误，无法使用密码登录")
			return
		}
		userInfo, err = i.authByPartner(params)
	} else {
		userInfo, err = i.authByKey(c, params)
	}
	if err != nil {
		util.JSON(c, "error", err.Error())
		return
	}
	token, err := util.GenerateToken(*userInfo)

	if err != nil {
		util.JSON(c, "error", "无效的参数", nil)
		return
	}
	util.JSON(c, "ok", "success", gin.H{
		"token": token,
	})
}

func (i *Auth) authByKey(c *gin.Context, params *LoginParam) (*util.User, error) {
	api, err := cloudflare.New(params.UserApiKey, params.Email)
	if err != nil {
		return nil, err
	}
	_, err = api.UserDetails(c)
	if err != nil {
		return nil, errors.New("key无效")
	}
	//accounts, _, err := api.Accounts(c, cloudflare.PaginationOptions{
	//	Page:    1,
	//	PerPage: 1,
	//})
	//if err != nil || len(accounts) <= 0 {
	//	return nil, errors.New("key无效")
	//}
	userInfo := &util.User{}
	//userInfo.ID = accounts[0].ID
	userInfo.Email = params.Email
	userInfo.UserApiKey = params.UserApiKey
	userInfo.AuthType = "user_api_key"
	return userInfo, nil
}

func (i *Auth) authByPartner(params *LoginParam) (*util.User, error) {

	v8cdnPost := util.V8cdnPostForm("https://api.cloudflare.com/host-gw.html", url.Values{
		"act":              {"user_auth"},
		"host_key":         {pkg.Conf.Cloudflare.HostKey},
		"cloudflare_email": {params.Email},
		"cloudflare_pass":  {params.Password},
	})
	logrus.WithFields(logrus.Fields{
		"data": v8cdnPost,
	}).Debug()
	data := &CloudflareResponse{}
	userInfo := &util.User{}
	err := json.Unmarshal([]byte(v8cdnPost), &data)

	if err != nil {
		return nil, errors.New("请求失败")
	}
	if data.Result != "success" {
		message := "未知异常"
		if data.Msg != nil {
			message = fmt.Sprintf("%v", data.Msg)
		}
		return nil, errors.New(message)
	}
	userInfo.Email = data.Response.CloudflareEmail
	userInfo.UserKey = data.Response.UserKey
	userInfo.UserApiKey = data.Response.UserApiKey
	userInfo.AuthType = "partner"
	return userInfo, nil
}
