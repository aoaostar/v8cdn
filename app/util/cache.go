package util

import (
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

func GetCacheDrive(username string) *cache.Cache {
	Cache := pkg.Cache[username]
	if Cache == nil {
		Cache = cache.New(24*time.Hour, 10*time.Second)
		pkg.Cache[username] = Cache
	}
	return Cache
}

func GetCache(c *gin.Context, key string) (interface{}, bool) {

	user := c.MustGet("user").(*User)
	Cache := GetCacheDrive(user.Email)
	key = c.Request.Method + "_" + c.FullPath() + "_" + c.ClientIP() + "_" + user.Email + "_" + key
	data, b := Cache.Get(key)
	if b {
		c.Header("v8cdn-cache", "HIT")
		return data, true
	}

	c.Header("v8cdn-cache", "MISS")
	return nil, false
}
func SetCache(c *gin.Context, key string, data interface{}) (interface{}, bool) {

	user := c.MustGet("user").(*User)
	Cache := GetCacheDrive(user.Email)
	key = c.Request.Method + "_" + c.FullPath() + "_" + c.ClientIP() + "_" + user.Email + "_" + key
	Cache.Set(key, data, 1*time.Minute)

	return nil, false
}
