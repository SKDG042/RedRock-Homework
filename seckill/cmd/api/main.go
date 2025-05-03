package main

import(
	"log"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/viper"

	"Redrock/seckill/internal/api/config"
	"Redrock/seckill/internal/api/client"
	"Redrock/seckill/internal/api/router"
	"Redrock/seckill/internal/pkg/redis"
)

func main(){

	// 读取配置
	viper.SetConfigName("api")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/api/config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil{
		log.Fatalf("读取API配置文件失败：%v", err)
	}

	var config config.Config
	err = viper.Unmarshal(&config)
	if err != nil{
		log.Fatalf("解析API配置文件失败：%v", err)
	}

	// 初始化Redis
	if err := redis.InitRedis(&config.Redis); err != nil{
		log.Fatalf("初始化Redis失败：%v", err)
	}
	defer redis.CloseRedis()

	// 初始化服务客户端
	clients, err := client.NewRPCClients(&config)
	if err != nil{
		log.Fatalf("初始化服务客户端失败：%v", err)
	}

	// 启动Hertz服务器
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)),
	)
	
	router.SetupRouter(h, clients)

	log.Printf("Hertz服务器启动成功，监听地址：%s:%d", config.Server.Host, config.Server.Port)
	h.Spin()
}
