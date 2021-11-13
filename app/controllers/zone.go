package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Zone struct {
}

// List 获取域名列表
func (z *Zone) List(c *gin.Context) {

	//获取缓存
	if zones, b := util.GetCache(c, "zones"); b {
		util.JSON(c, "ok", "success", zones)
		return
	}
	//缓存结束
	cf := c.MustGet("cloudflare").(*cloudflare.API)
	zones, err := cf.ListZones(c)

	log.WithFields(log.Fields{
		"data": zones,
	}).Debug()
	if err != nil {
		log.Debug(err)
		util.JSON(c, "error", err.Error(), err)
		return
	}
	//缓存结果
	util.SetCache(c, "zones", zones)
	util.JSON(c, "ok", "success", zones)
}

// Get 显示域名详情
func (z *Zone) Get(c *gin.Context) {

	var params ZoneParam
	if err := c.ShouldBind(&params); err != nil {
		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if params.ZoneId == "" {
		if ZoneId, err := cf.ZoneIDByName(params.ZoneName); err != nil {
			params.ZoneId = ZoneId
			log.Debug(params.ZoneId)
			util.JSON(c,
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			)
			return
		}
	}
	details, err := cf.ZoneDetails(c, params.ZoneId)
	log.WithFields(log.Fields{
		"data": details,
	}).Debug()
	if err != nil {
		util.JSON(c,
			"error",
			err.Error(), nil,
		)
		return
	}
	util.JSON(c, "ok", "success", details)
}

// Create 添加域名
func (z *Zone) Create(c *gin.Context) {

	type param struct {
		ZoneName string `json:"zone_name" valid:"required,url"`
	}
	var postData param

	if err := c.ShouldBind(&postData); err != nil {
		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	user := c.MustGet("user").(*util.User)
	if user.AuthType == "user_api_key" || user.UserKey == "" {
		util.JSON(c, "error", "user_api_key接入方式不支持添加域名")
		return
	}
	v8cdnPost := util.V8cdnPostForm("https://api.cloudflare.com/host-gw.html", url.Values{
		"act":        {"zone_set"},
		"host_key":   {pkg.Conf.Cloudflare.HostKey},
		"user_key":   {user.UserKey},
		"zone_name":  {postData.ZoneName},
		"resolve_to": {pkg.Conf.Cloudflare.DefaultRecord},
		"subdomains": {"@,www"},
	})
	log.WithFields(log.Fields{
		"data": v8cdnPost,
	}).Debug()
	var data gin.H
	err := json.Unmarshal([]byte(v8cdnPost), &data)

	if err != nil {
		util.JSON(c, "error", err.Error())
		return
	}
	if data["result"] != "success" {
		message := "未知异常"
		if data["msg"] != nil {
			message = fmt.Sprintf("%v", data["msg"])
		}
		util.JSON(c, "error", message, nil)
		return

	}
	if err != nil {
		util.JSON(c, "error", err.Error(), err)
		return
	}
	util.JSON(c, "ok", "success", data)

}

// Delete 删除域名
func (z *Zone) Delete(c *gin.Context) {

	var postData ZoneParam
	if err := c.ShouldBind(&postData); err != nil {
		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)
	if postData.ZoneId == "" {
		if ZoneId, err := cf.ZoneIDByName(postData.ZoneName); err != nil {
			postData.ZoneId = ZoneId
			log.Debug(postData.ZoneId)
			util.JSON(c,
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			)
			return
		}
	}
	data, err := cf.DeleteZone(c, postData.ZoneId)
	log.WithFields(log.Fields{
		"data": data,
	}).Debug()
	if err != nil {
		log.Debug(err)
		util.JSON(c, "error", err.Error(), err)
		return
	}

	util.JSON(c, "ok", "success", data)
}

// Update 更新域名
func (z *Zone) Update(c *gin.Context) {
}

func (z *Zone) createZoneByNs(c *gin.Context,ZoneName string) error {

	cf := c.MustGet("cloudflare").(*cloudflare.API)
	user := c.MustGet("user").(*util.User)
	account := cloudflare.Account{
		ID:       user.ID,
	}
	_, err := cf.CreateZone(c, ZoneName, false, account, "full")
	if err != nil {
		return err
	}
	return nil

}