package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/spf13/viper"

	"Redrock/seckill/internal/activity/config"
	"Redrock/seckill/internal/activity/service"
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/redis"
	activity "Redrock/seckill/kitex_gen/activity/internalactivityservice"
	"Redrock/seckill/internal/pkg/models"
)

func main(){

	// 读取配置
	viper.SetConfigName("activity")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/activity/config")
	viper.AutomaticEnv() // 自动读取环境变量，方便后续调试

	if err := viper.ReadInConfig(); err != nil{
		log.Fatalf("读取活动配置文件失败：%v", err)
	}

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil{
		log.Fatalf("解析活动配置文件失败：%v", err)
	}

	// 连接数据库
	if err := database.InitDB(&config.Database); err != nil{
		log.Fatalf("初始化连接数据库失败：%v", err)
	}
	defer database.CloseDB()

	// 迁移表结构
	if err := database.MigrateDB(&models.User{},&models.Activity{}, &models.Product{}, &models.Order{}); err != nil{
		log.Fatalf("迁移表结构失败：%v", err)
	}

	// 连接Redis
	if err := redis.InitRedis(&config.Redis); err != nil{
		log.Fatalf("初始化连接Redis失败：%v", err)
	}
	defer redis.CloseRedis()

	// 启动kitex服务
	activityImpl := service.NewInternalActivityServiceImpl()

	// 将字符串转化为TCP地址
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
	if err != nil{
		log.Fatalf("解析TCP地址失败：%v", err)
	}

	svr := activity.NewServer(
		activityImpl,
		server.WithServiceAddr(address),
		server.WithServerBasicInfo(nil),
	)

	if err := svr.Run(); err != nil{
		log.Fatalf("活动服务启动失败：%v", err)
	}
	
	log.Printf("活动服务启动成功，地址为：%s:%d", config.Server.Host, config.Server.Port)
}
