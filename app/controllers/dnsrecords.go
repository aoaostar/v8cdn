package controllers

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Dnsrecords struct {
}

// List 获取dns解析记录列表
func (d *Dnsrecords) List(c *gin.Context) {

	var postData ZoneParam
	if err := c.ShouldBind(&postData); err != nil {

		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}

	cf := c.MustGet("cloudflare").(*cloudflare.API)
	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			util.JSON(c,
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			)
			return
		}
	}
	//获取缓存
	if data, b := util.GetCache(c, postData.ZoneId); b {
		util.JSON(c,"ok", "success", data)
		return
	}
	data, err := cf.DNSRecords(c, postData.ZoneId, cloudflare.DNSRecord{})
	log.WithFields(log.Fields{
		"data": data,
	}).Debug()
	if err != nil {
		util.JSON(c,
			"error",
			fmt.Sprintf("请求失败：%v", err), err,
		)
		return
	}
	//缓存结果
	util.SetCache(c, postData.ZoneId, data)
	util.JSON(c,"ok", "success", data)
}

// Get 显示dns解析记录详情
func (d *Dnsrecords) Get(c *gin.Context) {

	type param struct {
		ZoneParam
		RecordId string `form:"record_id" json:"record_id" valid:"required,len=32"`
	}
	var postData param
	if err := c.ShouldBind(&postData); err != nil {

		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			util.JSON(c,
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			)
			return
		}
	}
	//获取缓存
	if data, b := util.GetCache(c, postData.ZoneId+"_"+postData.RecordId); b {
		util.JSON(c,"ok", "success", data)
		return
	}
	data, err := cf.DNSRecord(c, postData.ZoneId, postData.RecordId)

	log.WithFields(log.Fields{
		"data": data,
	}).Debug()
	if err != nil {
		util.JSON(c,
			"error",
			fmt.Sprintf("请求失败：%v", err), err,
		)
		return
	}
	//缓存结果
	util.SetCache(c, postData.ZoneId+"_"+postData.RecordId, data)
	util.JSON(c,"ok", "success", data)
}

// Create 添加dns解析记录
func (d *Dnsrecords) Create(c *gin.Context) {
	type param struct {
		ZoneParam
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
	if err := c.ShouldBind(&postData); err != nil {

		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}

	if postData.ZoneId == "" && postData.ZoneName == "" {
		util.JSON(c,
			"error",
			"The zone_id or zone_name is empty", nil,
		)
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
	record, err := cf.CreateDNSRecord(c, postData.ZoneId, cloudflare.DNSRecord{
		Type:     postData.DNSrecord.Type,
		Name:     postData.DNSrecord.Name,
		Content:  postData.DNSrecord.Content,
		TTL:      postData.DNSrecord.TTL,
		Priority: postData.DNSrecord.Priority,
		Proxied:  postData.DNSrecord.Proxied,
	})

	log.WithFields(log.Fields{
		"data": record,
	}).Debug()
	if err != nil {
		util.JSON(c,
			"error",
			fmt.Sprintf("添加记录失败：%v", err), err,
		)
		return
	}
	util.JSON(c,
		"ok",
		"success", nil,
	)
}

// Delete 删除dns解析记录
func (d *Dnsrecords) Delete(c *gin.Context) {

	type param struct {
		ZoneParam
		RecordId string `form:"record_id" json:"record_id" valid:"required,len=32"`
	}
	var postData param

	if err := c.ShouldBind(&postData); err != nil {

		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}
	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			util.JSON(c,
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			)
			return
		}
	}
	err := cf.DeleteDNSRecord(c, postData.ZoneId, postData.RecordId)

	log.WithFields(log.Fields{
		"data": err,
	}).Debug()
	if err != nil {
		util.JSON(c,
			"error",
			fmt.Sprintf("请求失败：%v", err), err,
		)
		return
	}
	util.JSON(c,"ok", "success", err)
}

// Update 更新dns解析记录
func (d *Dnsrecords) Update(c *gin.Context) {
	type param struct {
		ZoneParam
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
	if err := c.ShouldBind(&postData); err != nil {

		validateError, _ := util.FomateValidateError(err)
		util.JSON(c, "error", validateError)
		return
	}

	cf := c.MustGet("cloudflare").(*cloudflare.API)

	if postData.ZoneId == "" {
		ZoneId, err := cf.ZoneIDByName(postData.ZoneName)
		postData.ZoneId = ZoneId
		if err != nil {
			log.Debug(postData.ZoneId)
			util.JSON(c,
				"error",
				fmt.Sprintf("zoneId获取失败：%v", err), err,
			)
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
	log.WithFields(log.Fields{
		"data": err,
	}).Debug()
	if err != nil {
		util.JSON(c,
			"error",
			fmt.Sprintf("更新记录失败：%v", err), err,
		)
		return
	}
	util.JSON(c,
		"ok",
		"success", err,
	)
}
