package bootstrap

import (
	"github.com/aoaostar/v8cdn_panel/app/controllers"
	"github.com/aoaostar/v8cdn_panel/app/middleware"
	"github.com/aoaostar/v8cdn_panel/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(engine *gin.Engine) {
	InitApiRouter(engine)

}
func InitApiRouter(engine *gin.Engine) {
	engine.Use(middleware.CorsMiddleware)
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.FullPath())

	})
	auth := engine.Group("/auth", middleware.RateLimitMiddleware)
	{
		authController := new(controllers.Auth)
		auth.POST("/login", authController.Login)
	}

	api := engine.Group("/api", middleware.JWTAuthMiddleware)
	{
		ZoneController := new(controllers.Zone)
		DnsrecordsController := new(controllers.Dnsrecords)
		CacheController := new(controllers.Cache)
		SettingsController := new(controllers.Setting)

		api.Use(middleware.RateLimitMiddleware)
		api.Use(middleware.FlushCacheMiddleware)
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, util.Msg("ok", "success", nil))
		})
		api.GET("/zones", ZoneController.List)
		api.GET("/zone", ZoneController.Get)
		api.POST("/zone", ZoneController.Create)
		api.DELETE("/zone", ZoneController.Delete)
		api.PUT("/zone", ZoneController.Update)

		api.GET("/dnsrecords", DnsrecordsController.List)
		api.GET("/dnsrecord", DnsrecordsController.Get)
		api.POST("/dnsrecord", DnsrecordsController.Create)
		api.DELETE("/dnsrecord", DnsrecordsController.Delete)
		api.PUT("/dnsrecord", DnsrecordsController.Update)

		api.POST("/cache", CacheController.Clear)

		api.GET("/settings", SettingsController.List)
		api.GET("/setting", SettingsController.Get)
		api.PUT("/setting", SettingsController.Put)
		api.PATCH("/setting", SettingsController.Patch)
	}
}
