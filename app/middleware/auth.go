package middleware

import (
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, util.Msg("error", "请求头中Authorization为空", nil))
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, util.Msg("error", "请求头中auth格式有误", nil))
		c.Abort()
		return
	}
	// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
	mc, err := util.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.Msg("error", "无效的Token", nil))
		c.Abort()
		return
	}
	// 将当前请求的username信息保存到请求的上下文c上
	c.Set("username", mc.Username)
	c.Set("user_key", mc.UserKey)
	c.Set("user_api_key", mc.UserApiKey)
	api, err := cloudflare.New(c.MustGet("user_api_key").(string), c.MustGet("username").(string))
	c.Set("cloudflare", api)
	c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
}
