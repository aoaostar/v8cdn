package bootstrap

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/config"
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
	engine.Use(static.Serve("/", static.LocalFile(config.Conf.Static, false)))
	engine.NoRoute(func(c *gin.Context) {
		c.File(config.Conf.Static + "/index.html")
	})
	//初始化路由
	InitRouter(engine)
	//初始化数据库
	//初始化缓存
	err := engine.Run(config.Conf.Listen)
	if err != nil {
		log.Error("Gin启动失败：%+v\n", err)
		return
	}
	fmt.Println("启动成功")
	log.Info("启动成功")

}
