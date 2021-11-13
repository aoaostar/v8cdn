package bootstrap

import (
	"github.com/aoaostar/v8cdn_panel/pkg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func InitLog() {

	log.SetFormatter(&log.JSONFormatter{
		//ForceColors:               true,
		//EnvironmentOverrideColors: true,
		TimestampFormat: "2006-01-02 15:04:05",
		//FullTimestamp:             true,
	})

	logFile := "logs/v8cdn.log"

	logWriter := &lumberjack.Logger{
		Filename:   logFile, //日志文件位置
		MaxSize:    5,       // 单文件最大容量,单位是MB
		MaxBackups: 3,       // 最大保留过期文件个数
		MaxAge:     7,       // 保留过期文件的最大时间间隔,单位是天
		Compress:   false,   // 是否需要压缩滚动日志, 使用的 gzip 压缩
		LocalTime:  true,
	}

	log.SetOutput(io.MultiWriter(logWriter, os.Stdout))

	if pkg.Conf.Debug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}
}
