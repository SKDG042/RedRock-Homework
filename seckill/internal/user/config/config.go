package config

import (
	"Redrock/seckill/internal/pkg/database"
)

// Config 定义了用户服务所需的配置
type Config struct {
	Server   ServerConfig      			`mapstructure:"server"`
	Database database.DatabaseConfig 	`mapstructure:"database"`
}

// ServerConfig 定义了Kitex服务器的配置
type ServerConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Host        string `mapstructure:"host"` 
	Port        int    `mapstructure:"port"` 
	LogLevel    string `mapstructure:"log_level"`
}
