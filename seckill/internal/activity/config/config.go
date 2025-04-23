package config

import(
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/redis"
)

// kitex服务器配置
type ServerConfig struct{
	ServiceName string 	`mapstructure:"service_name"`
	Host 		string 	`mapstructure:"host"`
	Port 		int 	`mapstructure:"port"`
	LogLevel 	string 	`mapstructure:"log_level"` // 日志等级从低到高: debug, info, warn, error
}
type Config struct{
	Server 		ServerConfig 				`mapstructure:"server"` 
	Database 	database.DatabaseConfig 	`mapstructure:"database"`
	Redis 		redis.RedisConfig 			`mapstructure:"redis"`
}
