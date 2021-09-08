package middleware

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/aoaostar/v8cdn_panel/config"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/patrickmn/go-cache"
	"net/http"
)

func RateLimitMiddleware(c *gin.Context) {
	if !config.Conf.RateLimit.Enabled {
		c.Next()
		return
	}
	key := c.Request.Method + "_" + c.FullPath() + "_" + c.ClientIP()

	rateLimitBucket, b := config.RateLimitCache.Get(key)
	if !b {
		rateLimitBucket = config.RateLimitBucket
	}

	c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%v", config.Conf.RateLimit.Capacity))
	c.Header("X-Rate-Limit-Duration", fmt.Sprintf("%v", config.Conf.RateLimit.FillInterval))
	c.Header("X-Rate-Limit-Available", fmt.Sprintf("%v", rateLimitBucket.(*ratelimit.Bucket).Available()))
	c.Header("X-Rate-Limit-Request-Key", key)

	if rateLimitBucket.(*ratelimit.Bucket).TakeAvailable(1) < 1 {
		c.JSON(http.StatusOK, util.Msg("error", "You have reached maximum request limit.", nil))
		c.Abort()
		return
	}
	config.RateLimitCache.Set(key, rateLimitBucket, cache.NoExpiration)
	c.Next()
}
