package controllers

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Setting struct {
}

func (s *Setting) List(c *gin.Context) {

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
	zoneSettings, err := cf.ZoneSettings(c, params.ZoneId)
	log.WithFields(log.Fields{
		"zoneSettings": zoneSettings,
	}).Debug()
	if err != nil {
		util.JSON(c, "error", err.Error(), err)
		return
	}
	util.JSON(c, "ok", "success", zoneSettings)
}
func (s *Setting) Get(c *gin.Context) {

	var params struct {
		ZoneParam
		SettingName string `form:"setting_name" json:"setting_name" binding:"required"`
	}
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
	zoneSingleSetting, err := cf.ZoneSingleSetting(c, params.ZoneId, params.SettingName)
	log.WithFields(log.Fields{
		"zoneSingleSetting": zoneSingleSetting,
	}).Debug()
	if err != nil {
		util.JSON(c, "error", err.Error(), err)
		return
	}
	util.JSON(c, "ok", "success", zoneSingleSetting)

}
func (s *Setting) Put(c *gin.Context) {

	var params struct {
		ZoneParam
		Settings []cloudflare.ZoneSetting `form:"settings" json:"settings" binding:"required"`
	}
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
	settingResponse, err := cf.UpdateZoneSettings(c, params.ZoneId, params.Settings)
	log.WithFields(log.Fields{
		"settingResponse": settingResponse,
	}).Debug()
	if err != nil {
		util.JSON(c, "error", err.Error(), err)
		return
	}
	util.JSON(c, "ok", "success", settingResponse)

}
func (s *Setting) Patch(c *gin.Context) {

	var params struct {
		ZoneParam
		SettingName string                 `form:"setting_name" json:"setting_name" binding:"required"`
		Setting     cloudflare.ZoneSetting `form:"setting" json:"setting" binding:"required"`
	}
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
	settingResponse, err := cf.UpdateZoneSingleSetting(c, params.ZoneId, params.SettingName, params.Setting)
	log.WithFields(log.Fields{
		"settingResponse": settingResponse,
	}).Debug()
	if err != nil {
		util.JSON(c, "error", err.Error(), err)
		return
	}
	util.JSON(c, "ok", "success", settingResponse)

}
