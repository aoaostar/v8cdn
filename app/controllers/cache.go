package controllers

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Cache struct {
}

type clearParam struct {
	ZoneParam
	Everything bool     `form:"everything" json:"everything" binding:"omitempty,required_without=Files"`
	Files      []string `form:"files" json:"files" binding:"omitempty"`
}

func (i *Cache) Clear(c *gin.Context) {
	api := c.MustGet("cloudflare").(*cloudflare.API)

	var params clearParam

	if err := c.ShouldBind(&params); err != nil {
		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)
	if params.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(params.ZoneName)
		params.ZoneId = ZoneId
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ZoneId": ZoneId,
				"err":    err,
			}).Debug()
			util.JSON(c, "error", fmt.Sprintf("zoneId获取失败：%v", err), err)
			return
		}
	}
	if len(params.Files) > 0 {
		params.Everything = false
	}
	pcr := cloudflare.PurgeCacheRequest{
		Everything: params.Everything,
		Files:      params.Files,
	}
	res, err := api.PurgeCache(c, params.ZoneId, pcr)
	logrus.WithFields(logrus.Fields{
		"data": res,
	}).Debug()
	if err != nil {
		util.JSON(c, "error", err.Error())
		return
	}
	util.JSON(c, "ok", "success")
	return

}
