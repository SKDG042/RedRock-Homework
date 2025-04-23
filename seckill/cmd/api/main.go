package main

import(
	"log"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/viper"

	"Redrock/seckill/internal/api/config"
)

func main(){

	// 读取配置
	viper.SetConfigName("activity")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/activity/config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil{
		log.Fatalf("读取活动配置文件失败：%v", err)
	}

	var config config.Config
	err = viper.Unmarshal(&config)
	if err != nil{
		log.Fatalf("解析活动配置文件失败：%v", err)
	}

	// 启动Hertz服务器
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)),
	)
	// TODO:启动kitex服务
	log.Printf("Hertz服务器启动成功，监听地址：%s:%d", config.Server.Host, config.Server.Port)
	h.Spin()
}
