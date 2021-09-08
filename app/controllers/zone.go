package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/aoaostar/v8cdn_panel/config"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type param struct {
	ZoneName string `form:"zone_name" json:"zone_name" binding:"omitempty,required_without=ZoneId,hostname"`
	ZoneId   string `form:"zone_id" json:"zone_id" binding:"omitempty,required_without=ZoneName,len=32"`
}

// ShowZones 获取域名列表
func ShowZones(c *gin.Context) {

	//获取缓存
	if zones, b := util.GetCache(c, "zones"); b {
		c.JSON(http.StatusOK, util.Msg("ok", "success", zones))
		return
	}
	//缓存结束
	cf := c.MustGet("cloudflare").(*cloudflare.API)
	zones, err := cf.ListZones(c)

	log.Debug(zones)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusOK, util.Msg("error", err.Error(), err))
		return
	}
	//缓存结果
	util.SetCache(c, "zones", zones)
	c.JSON(http.StatusOK, util.Msg("ok", "success", zones))
}

// DetailZone 显示域名详情
func DetailZone(c *gin.Context) {

	var postData param
	if err := c.ShouldBindQuery(&postData); err != nil {
		util.PrintError(err,c)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		if ZoneId, err := cf.ZoneIDByName(postData.ZoneName); err != nil {
			postData.ZoneId = ZoneId
			log.Debug(postData.ZoneId)
			c.JSON(http.StatusOK, util.Msg(
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			))
			return
		}
	}
	details, err := cf.ZoneDetails(c, postData.ZoneId)
	log.Debug(details)
	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			err.Error(), nil,
		))
		return
	}
	c.JSON(http.StatusOK, util.Msg("ok", "success", details))
}

// CreateZone 添加域名
func CreateZone(c *gin.Context) {

	type param struct {
		ZoneName string `json:"zone_name" valid:"required,url"`
	}
	var postData param

	if err := c.BindJSON(&postData); err != nil {
		util.PrintError(err,c)
		return
	}
	v8cdnPost := util.V8cdnPostForm("https://api.cloudflare.com/host-gw.html", url.Values{
		"act":        {"zone_set"},
		"host_key":   {config.Conf.Cloudflare.HostKey},
		"user_key":   {c.MustGet("user_key").(string)},
		"zone_name":  {postData.ZoneName},
		"resolve_to": {"v8cdn.cc"},
		"subdomains": {"@,www"},
	})
	log.Debug(v8cdnPost)
	var data gin.H
	err := json.Unmarshal([]byte(v8cdnPost), &data)

	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			err.Error(), nil,
		))
		return
	}
	if data["result"] != "success" {
		message := "未知异常"
		if data["msg"] != nil {
			message = fmt.Sprintf("%v", data["msg"])
		}
		c.JSON(http.StatusOK, util.Msg("error", message, nil))
		return

	}
	if err != nil {
		c.JSON(http.StatusOK, util.Msg("error", err.Error(), err))
		return
	}
	c.JSON(http.StatusOK, util.Msg("ok", "success", data))

}

// DeleteZone 删除域名
func DeleteZone(c *gin.Context) {

	var postData param
	if err := c.ShouldBindQuery(&postData); err != nil {
		util.PrintError(err,c)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)
	if postData.ZoneId == "" {
		if ZoneId, err := cf.ZoneIDByName(postData.ZoneName); err != nil {
			postData.ZoneId = ZoneId
			log.Debug(postData.ZoneId)
			c.JSON(http.StatusOK, util.Msg(
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			))
			return
		}
	}
	data, err := cf.DeleteZone(c, postData.ZoneId)
	log.Debug(data)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusOK, util.Msg("error", err.Error(), err))
		return
	}

	c.JSON(http.StatusOK, util.Msg("ok", "success", data))
}

// UpdateZone 更新域名
func UpdateZone(c *gin.Context) {
}
