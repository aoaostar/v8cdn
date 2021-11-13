package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Msg(status string, message string, data ...interface{}) gin.H {

	var data2 interface{}
	if len(data) == 0 {
		data2 = gin.H{}
	} else {
		if data[0] == nil {
			data2 = gin.H{}
		} else {
			data2 = data[0]
		}
	}
	return gin.H{
		"status":  status,
		"message": message,
		"data":    data2,
	}

}
func JSON(c *gin.Context, status string, message string, data ...interface{}) {

	var data2 interface{}
	if len(data) == 0 {
		data2 = gin.H{}
	} else {
		if data[0] == nil {
			data2 = gin.H{}
		} else {
			data2 = data[0]
		}
	}
	c.JSON(http.StatusOK, Msg(status, message, data2))
}