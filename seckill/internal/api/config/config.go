package config

import (
	"Redrock/seckill/internal/pkg/redis"
)

type Config struct{
	Server		ServerConfig		`mapstructure:"server"`
	UserRPC		ClientConfig		`mapstructure:"user_rpc"`
	ActivityRPC	ClientConfig		`mapstructure:"activity_rpc"`
	OrderRPC	ClientConfig		`mapstructure:"order_rpc"`
	Redis		redis.RedisConfig	`mapstructure:"redis"`
}

// 这里为Hertz服务器的配置
type ServerConfig struct{
	Host 		string 	`mapstructure:"host"`
	Port 		int 	`mapstructure:"port"`
	LogLevel 	string 	`mapstructure:"log_level"`
	RateLimit 	int 	`mapstructure:"rate_limit"` // 限制每秒请求数
}

// 这里是kitex client的配置
type ClientConfig struct{
	ServiceName string `mapstructure:"service_name"`
	TargetHost  string `mapstructure:"target_host"`
	TargetPort  int    `mapstructure:"target_port"`
	Timeout     int    `mapstructure:"timeout"`
}
