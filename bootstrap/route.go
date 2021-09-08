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
		c.JSON(http.StatusOK,c.FullPath())

	})
	auth := engine.Group("/auth", middleware.RateLimitMiddleware)
	{

		auth.POST("/login", controllers.Login)
		auth.POST("/parse", controllers.Parse)
	}

	api := engine.Group("/api", middleware.JWTAuthMiddleware)
	{
		api.Use(middleware.RateLimitMiddleware)
		api.Use(middleware.FlushCacheMiddleware)
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, util.Msg("ok", "success", nil))
		})
		api.GET("/zones", controllers.ShowZones)
		api.GET("/zone", controllers.DetailZone)
		api.POST("/zone", controllers.CreateZone)
		api.DELETE("/zone", controllers.DeleteZone)
		api.PUT("/zone", controllers.UpdateZone)

		api.GET("/dnsrecords", controllers.ShowDNSRecords)
		api.GET("/dnsrecord", controllers.DetailDNSRecord)
		api.POST("/dnsrecord", controllers.CreateDNSRecord)
		api.DELETE("/dnsrecord", controllers.DeleteDNSRecord)
		api.PUT("/dnsrecord", controllers.UpdateDNSRecord)

		api.GET("/AccountDetails", controllers.AccountDetails)
		api.GET("/UserDetails", controllers.UserDetails)
		api.GET("/CreateZone", controllers.CreateZone)
	}
}
