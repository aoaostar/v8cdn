package util

import "github.com/gin-gonic/gin"

func Msg(status string, message string, data interface{}) gin.H {
	if data == nil {
		data = []string{}
	}
	return gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	}
}
