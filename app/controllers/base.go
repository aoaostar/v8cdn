package controllers

type ZoneParam struct {
	ZoneId   string `form:"zone_id" json:"zone_id" binding:"required_without=ZoneName,len=32"`
	ZoneName string `form:"zone_name" json:"zone_name" binding:"omitempty,hostname"`
}

