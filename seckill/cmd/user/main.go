package main

import (
	"log"
	"net"
	"fmt"

	"github.com/cloudwego/kitex/server"
	"github.com/spf13/viper"

	userService "Redrock/seckill/kitex_gen/user/userservice"
	"Redrock/seckill/internal/user/config"
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/models"
	"Redrock/seckill/internal/user/service"
	
)

func main() {
	// 加载配置
	viper.SetConfigName("user")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/user/config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("读取用户配置文件失败: %v", err)
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("解析用户配置文件失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.CloseDB()

	// 自动迁移数据库表结构
	if err := database.MigrateDB(&models.User{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建服务实现实例
	userImpl := service.NewUserServiceImpl()

	// 创建Kitex服务器
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatalf("解析TCP地址失败: %v", err)
	}

	svr := userService.NewServer(
		userImpl, 
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(nil),
	)

	// 启动Kitex服务器
	if err := svr.Run(); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}

	log.Printf("用户服务启动成功，地址为：%s:%d", cfg.Server.Host, cfg.Server.Port)
}