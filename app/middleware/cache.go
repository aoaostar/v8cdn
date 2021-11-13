package middleware

import (
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FlushCacheMiddleware(c *gin.Context) {
	methods := [...]string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
	user := c.MustGet("user").(*util.User)
	for _, v := range methods {
		if c.Request.Method == v {
			util.GetCacheDrive(user.Email).Flush()
			break
		}

	}
	c.Next()
}
