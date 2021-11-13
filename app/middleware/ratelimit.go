package middleware

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/patrickmn/go-cache"
	"net/http"
)

func RateLimitMiddleware(c *gin.Context) {
	if !pkg.Conf.RateLimit.Enabled {
		c.Next()
		return
	}
	key := c.Request.Method + "_" + c.FullPath() + "_" + c.ClientIP()

	rateLimitBucket, b := pkg.RateLimitCache.Get(key)
	if !b {
		rateLimitBucket = pkg.RateLimitBucket
	}

	c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%v", pkg.Conf.RateLimit.Capacity))
	c.Header("X-Rate-Limit-Duration", fmt.Sprintf("%v", pkg.Conf.RateLimit.FillInterval))
	c.Header("X-Rate-Limit-Available", fmt.Sprintf("%v", rateLimitBucket.(*ratelimit.Bucket).Available()))
	c.Header("X-Rate-Limit-Request-UserApiKey", key)

	if rateLimitBucket.(*ratelimit.Bucket).TakeAvailable(1) < 1 {
		c.JSON(http.StatusOK, util.Msg("error", "You have reached maximum request limit.", nil))
		c.Abort()
		return
	}
	pkg.RateLimitCache.Set(key, rateLimitBucket, cache.NoExpiration)
	c.Next()
}
