package bootstrap

import (
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/static"
)

func Run() {
	InitConfig()
	InitLog()
	InitRateLimit()
	InitValidator()
	engine := gin.Default()
	if !pkg.Conf.Debug {
		//记录panic错误
		engine.Use(gin.RecoveryWithWriter(log.StandardLogger().Writer()))
	}
	engine.Use(static.Serve("/", static.LocalFile(pkg.Conf.Static, false)))
	engine.NoRoute(func(c *gin.Context) {
		c.File(pkg.Conf.Static + "/index.html")
	})
	//初始化路由
	InitRouter(engine)
	//初始化数据库
	//初始化缓存
	err := engine.Run(pkg.Conf.Listen)
	if err != nil {
		log.Error("Gin启动失败：%+v\n", err)
		return
	}

}
