package controllers

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)


// ShowDNSRecords 获取dns解析记录列表
func ShowDNSRecords(c *gin.Context) {

	type param struct {
		ZoneName string `form:"zone_name" json:"zone_name" binding:"omitempty,required_without=ZoneId,hostname"`
		ZoneId   string `form:"zone_id" json:"zone_id" binding:"omitempty,required_without=ZoneName,len=32"`
	}
	var postData param
	if err := c.ShouldBindQuery(&postData); err != nil {
		util.PrintError(err,c)
		return
	}

	cf := c.MustGet("cloudflare").(*cloudflare.API)
	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			c.JSON(http.StatusOK, util.Msg(
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			))
			return
		}
	}
	//获取缓存
	if data, b := util.GetCache(c, postData.ZoneId); b {
		c.JSON(http.StatusOK, util.Msg("ok", "success", data))
		return
	}
	data, err := cf.DNSRecords(c, postData.ZoneId, cloudflare.DNSRecord{})
	log.Debug(data)
	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			fmt.Sprintf("请求失败：%v", err), err,
		))
		return
	}
	//缓存结果
	util.SetCache(c, postData.ZoneId, data)
	c.JSON(http.StatusOK, util.Msg("ok", "success", data))
}

// DetailDNSRecord 显示dns解析记录详情
func DetailDNSRecord(c *gin.Context) {

	type param struct {
		ZoneName string `form:"zone_name" json:"zone_name" valid:"omitempty,required_without=ZoneId,hostname"`
		ZoneId   string `form:"zone_id" json:"zone_id" valid:"omitempty,required_without=ZoneName,len=32"`
		RecordId string `form:"record_id" json:"record_id" valid:"required,len=32"`
	}
	var postData param
	if err := c.ShouldBindQuery(&postData); err != nil {
		util.PrintError(err,c)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			c.JSON(http.StatusOK, util.Msg(
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			))
			return
		}
	}
	//获取缓存
	if data, b := util.GetCache(c, postData.ZoneId+"_"+postData.RecordId); b {
		c.JSON(http.StatusOK, util.Msg("ok", "success", data))
		return
	}
	data, err := cf.DNSRecord(c, postData.ZoneId, postData.RecordId)
	log.Debug(data)
	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			fmt.Sprintf("请求失败：%v", err), err,
		))
		return
	}
	//缓存结果
	util.SetCache(c, postData.ZoneId+"_"+postData.RecordId, data)
	c.JSON(http.StatusOK, util.Msg("ok", "success", data))
}

// CreateDNSRecord 添加dns解析记录
func CreateDNSRecord(c *gin.Context) {
	type param struct {
		ZoneName  string `json:"zone_name" valid:"omitempty,required_without=ZoneId,hostname"`
		ZoneId    string `form:"zone_id" json:"zone_id" valid:"omitempty,required_without=ZoneName,len=32"`
		DNSrecord struct {
			Type     string  `json:"type" valid:"required,oneof(A AAAA CNAME HTTPS TXT SRV LOC MX NS SPF CERT DNSKEY DS NAPTR SMIMEA SSHFP SVCB TLSA URI)"`
			Name     string  `json:"name" valid:"required,min=1,max=255)"`
			Content  string  `json:"content" valid:"required"`
			TTL      int     `json:"ttl" valid:"required,numeric"`
			Priority *uint16 `json:"priority" valid:"required,min=0,max=65535"`
			Proxied  *bool   `json:"proxied" valid:"required,bool"`
		}
	}
	var postData param
	if err := c.BindJSON(&postData); err != nil {
		util.PrintError(err,c)
		return
	}

	if postData.ZoneId == "" && postData.ZoneName == "" {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			"The zone_id or zone_name is empty", nil,
		))
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
	record, err := cf.CreateDNSRecord(c, postData.ZoneId, cloudflare.DNSRecord{
		Type:     postData.DNSrecord.Type,
		Name:     postData.DNSrecord.Name,
		Content:  postData.DNSrecord.Content,
		TTL:      postData.DNSrecord.TTL,
		Priority: postData.DNSrecord.Priority,
		Proxied:  postData.DNSrecord.Proxied,
	})
	log.Debug(record)
	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			fmt.Sprintf("添加记录失败：%v", err), err,
		))
		return
	}
	c.JSON(http.StatusOK, util.Msg(
		"ok",
		"success", nil,
	))
}

// DeleteDNSRecord 删除dns解析记录
func DeleteDNSRecord(c *gin.Context) {

	type param struct {
		ZoneName string `json:"zone_name" valid:"omitempty,required_without=ZoneId,hostname"`
		ZoneId   string `form:"zone_id" json:"zone_id" valid:"omitempty,required_without=ZoneName,len=32"`
		RecordId string `form:"record_id" json:"record_id" valid:"required,len=32"`
	}
	var postData param

	if err := c.ShouldBindQuery(&postData); err != nil {
		util.PrintError(err,c)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			c.JSON(http.StatusOK, util.Msg(
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			))
			return
		}
	}
	err := cf.DeleteDNSRecord(c, postData.ZoneId, postData.RecordId)

	log.Debug(err)
	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			fmt.Sprintf("请求失败：%v", err), err,
		))
		return
	}
	c.JSON(http.StatusOK, util.Msg("ok", "success", err))
}

// UpdateDNSRecord 更新dns解析记录
func UpdateDNSRecord(c *gin.Context) {
	type param struct {
		ZoneName  string `json:"zone_name" valid:"omitempty,required_without=ZoneId,hostname"`
		ZoneId    string `form:"zone_id" json:"zone_id" valid:"omitempty,required_without=ZoneName,len=32"`
		RecordId  string `form:"record_id" json:"record_id" valid:"required,len=32"`
		DNSrecord struct {
			Type     string  `json:"type" valid:"omitempty,oneof(A AAAA CNAME HTTPS TXT SRV LOC MX NS SPF CERT DNSKEY DS NAPTR SMIMEA SSHFP SVCB TLSA URI)"`
			Name     string  `json:"name" valid:"omitempty,min=1,max=255)"`
			Content  string  `json:"content" valid:"omitempty,min=1"`
			TTL      int     `json:"ttl" valid:"omitempty,numeric"`
			Priority *uint16 `json:"priority" valid:"omitempty,min=0,max=65535"`
			Proxied  *bool   `json:"proxied" valid:"omitempty,bool"`
		}
	}
	var postData param
	if err := c.BindJSON(&postData); err != nil {
		util.PrintError(err,c)
		return
	}

	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			c.JSON(http.StatusOK, util.Msg(
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			))
			return
		}
	}
	err := cf.UpdateDNSRecord(c, postData.ZoneId, postData.RecordId, cloudflare.DNSRecord{
		Type:     postData.DNSrecord.Type,
		Name:     postData.DNSrecord.Name,
		Content:  postData.DNSrecord.Content,
		TTL:      postData.DNSrecord.TTL,
		Priority: postData.DNSrecord.Priority,
		Proxied:  postData.DNSrecord.Proxied,
	})
	log.Debug(err)
	if err != nil {
		c.JSON(http.StatusOK, util.Msg(
			"error",
			fmt.Sprintf("更新记录失败：%v", err), err,
		))
		return
	}
	c.JSON(http.StatusOK, util.Msg(
		"ok",
		"success", err,
	))
}
