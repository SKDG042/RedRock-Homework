package config

import (
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/mq"
	"Redrock/seckill/internal/pkg/redis"
)

type Config struct{
	Server		ServerConfig			`mapstructure:"server"`
	Database	database.DatabaseConfig	`mapstructure:"database"`
	Redis		redis.RedisConfig		`mapstructure:"redis"`
	MQ			mq.MQConfig				`mapstructure:"mq"`
	ActivityRPC	ActivityRPCConfig		`mapstructure:"activity_rpc"`
}

type ActivityRPCConfig struct{
	Host		string	`mapstructure:"host"`
	Port		int		`mapstructure:"port"`
	Timeout		int		`mapstructure:"timeout"`
}

type ServerConfig struct{
	ServiceName		string		`mapstructure:"service_name"`
	Host			string		`mapstructure:"host"`
	Port			int			`mapstructure:"port"`
	LogLevel		string		`mapstructure:"log_level"`
}
