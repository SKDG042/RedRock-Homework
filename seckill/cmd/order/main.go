package main

import (
	"log"

	"github.com/spf13/viper"

	"Redrock/seckill/internal/order/config"
	"Redrock/seckill/internal/order/service"
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/redis"
	order "Redrock/seckill/kitex_gen/order/orderservice"
)

func main(){

	// 读取配置
	viper.SetConfigName("order")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/order/config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil{
		log.Fatalf("读取订单配置文件失败：%v", err)
	}

	var config config.Config
	err = viper.Unmarshal(&config)
	if err != nil{
		log.Fatalf("解析订单配置文件失败：%v", err)
	}

	// 连接数据库
	if err := database.InitDB(&config.Database); err != nil{
		log.Fatalf("初始化连接数据库失败：%v", err)
	}
	defer database.CloseDB()

	// 连接Redis
	if err := redis.InitRedis(&config.Redis); err != nil{
		log.Fatalf("初始化连接Redis失败：%v", err)
	}
	defer redis.CloseRedis()

	// 连接MQ todo...

	// 启动kitex服务
	svr := order.NewServer(new(service.OrderServiceImpl))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}

	select{}
}
