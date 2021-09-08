package util

import (
	"github.com/aoaostar/v8cdn_panel/config"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

func GetCacheDrive(username string) *cache.Cache {
	Cache := config.Cache[username]
	if Cache == nil {
		Cache = cache.New(24*time.Hour, 10*time.Second)
		config.Cache[username] = Cache
	}
	return Cache
}

func GetCache(c *gin.Context, key string) (interface{}, bool) {

	username := c.MustGet("username").(string)
	Cache := GetCacheDrive(username)
	key = c.Request.Method + "_" + c.FullPath() + "_" + c.ClientIP() + "_" + username + "_" + key
	data, b := Cache.Get(key)
	if b {
		c.Header("v8cdn-cache", "HIT")
		return data, true
	}

	c.Header("v8cdn-cache", "MISS")
	return nil, false
}
func SetCache(c *gin.Context, key string, data interface{}) (interface{}, bool) {

	username := c.MustGet("username").(string)
	Cache := GetCacheDrive(username)
	key = c.Request.Method + "_" + c.FullPath() + "_" + c.ClientIP() + "_" + username + "_" + key
	Cache.Set(key, data, 1*time.Minute)

	return nil, false
}
