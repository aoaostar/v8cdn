package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
	"os"
	"time"
)

func InitLog() {

	log.SetFormatter(&log.JSONFormatter{
		//ForceColors:               true,
		//EnvironmentOverrideColors: true,
		TimestampFormat: "2006-01-02 15:04:05",
		//FullTimestamp:             true,
	})
	log_dir := "./logs"
	_, err := os.Stat(log_dir)
	if err != nil {
		if err := os.Mkdir(log_dir, os.ModePerm); err != nil {
			log.Debugf("日志目录创建失败：%s", err)
		}
	}
	if config.Get("debug") == true {

		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
		//log.SetReportCaller(true)

	}else{
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
		f, err := os.OpenFile(fmt.Sprintf("%s/%s.log", log_dir, time.Now().Format("2006-01-02")), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			log.Error("Failed to log to file")
		} else {

			log.SetOutput(f)
		}

	}
}
