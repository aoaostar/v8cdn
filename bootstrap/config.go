package bootstrap

import (
	"fmt"
	"github.com/aoaostar/v8cdn_panel/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil { // 处理读取配置文件的错误
		log.Panic(fmt.Errorf("读取配置出错: %s \n", err))
	}
	if err := viper.Unmarshal(&pkg.Conf); err != nil {
		log.Panic(fmt.Errorf("解析配置出错: %s \n", err))
	}

}
