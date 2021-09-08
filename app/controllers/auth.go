package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/gin-gonic/gin"
	config "github.com/spf13/viper"
	"net/http"
	"net/url"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
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

func Login(c *gin.Context) {
	//https://api.cloudflare.com/host-gw.html
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Msg("error", "无效的参数", nil))
		return
	}
	v8cdnPost := util.V8cdnPostForm("https://api.cloudflare.com/host-gw.html", url.Values{
		"act":              {"user_auth"},
		"host_key":         {config.GetString("cloudflare.host_key")},
		"cloudflare_email": {user.Username},
		"cloudflare_pass":  {user.Password},
	})
	var data CloudflareResponse
	err = json.Unmarshal([]byte(v8cdnPost), &data)

	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			"请求失败", nil,
		))
		return
	}
	if data.Result != "success" {
		message := "未知异常"
		if data.Msg != nil {
			message = fmt.Sprintf("%v", data.Msg)
		}
		c.JSON(http.StatusOK, util.Msg("error", message, nil))
		return

	}
	token, err := util.GenerateToken(data.Response.CloudflareEmail, data.Response.UserKey, data.Response.UserApiKey)

	if err != nil {

		c.JSON(http.StatusOK, util.Msg("error", "无效的参数", nil))
		return
	}
	c.JSON(http.StatusOK, util.Msg("ok", "success", gin.H{
		"token": token,
	}))

}
func Parse(c *gin.Context) {
	token := c.PostForm("token")
	c.JSON(http.StatusOK, gin.H{
		"data": token,
	})
	parseToken, _ := util.ParseToken(token)

	c.JSON(http.StatusOK, gin.H{
		"data": parseToken,
	})

}
