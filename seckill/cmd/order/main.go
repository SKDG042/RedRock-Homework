package main

import (
	"fmt"
	"net"
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
	"github.com/spf13/viper"

	"Redrock/seckill/internal/order/config"
	"Redrock/seckill/internal/order/data"
	"Redrock/seckill/internal/order/mq"
	"Redrock/seckill/internal/order/service"
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/models"
	"Redrock/seckill/internal/pkg/redis"
	activityClient "Redrock/seckill/kitex_gen/activity/activityservice"
	order "Redrock/seckill/kitex_gen/order/orderservice"
)

func main(){

	// 读取配置
	viper.SetConfigName("order")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/order/config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil{
		log.Fatalf("读取订单配置文件失败：%v", err)
	}

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil{
		log.Fatalf("解析订单配置文件失败：%v", err)
	}

	// 连接数据库
	if err := database.InitDB(&config.Database); err != nil{
		log.Fatalf("初始化连接数据库失败：%v", err)
	}
	defer database.CloseDB()

	// 自动迁移数据库表
	if err := database.MigrateDB(&models.User{}, &models.Product{}, &models.Activity{}, &models.Order{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 连接Redis
	if err := redis.InitRedis(&config.Redis); err != nil{
		log.Fatalf("初始化连接Redis失败：%v", err)
	}
	defer redis.CloseRedis()

	// 启动ActivityService的客户端
	activityServiceClient, err := activityClient.NewClient(
		"activity_service",
		client.WithHostPorts(fmt.Sprintf("%s:%d", config.ActivityRPC.Host, config.ActivityRPC.Port)),
		client.WithRPCTimeout(time.Duration(config.ActivityRPC.Timeout) * time.Millisecond),
	)

	if err != nil{
		log.Fatalf("连接ActivityService客户端失败：%v", err)
	}

	// 初始化订单消息生产者
	orderProducer, err := mq.NewOrderProducer(&config.MQ.RabbitMQ)
	if err != nil{
		log.Fatalf("初始化订单消息生产者失败：%v", err)
	}

	defer orderProducer.Close()

	// 初始化订单消费者
	orderConsumer, err := mq.NewOrderConsumer(&config.MQ.RabbitMQ, data.NewOrderData())
	if err != nil{
		log.Fatalf("初始化订单消费者失败：%v", err)
	}

	defer orderConsumer.Close()

	// 启动消费者
	if err := orderConsumer.StartConsume(); err != nil{
		log.Fatalf("启动消费者失败：%v", err)
	}

	// 启动kitex服务
	orderImpl := service.NewOrderServiceImpl(orderProducer, activityServiceClient)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
	if err != nil{
		log.Fatalf("解析TCP地址失败: %v", err)
	}

	svr := order.NewServer(
		orderImpl,
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(nil),
	)

	err = svr.Run()

	if err != nil {
		log.Fatalf("启动kitex服务器失败：%v", err)
	}

	log.Printf("订单服务启动成功，监听地址：%s:%d", config.Server.Host, config.Server.Port)
}
