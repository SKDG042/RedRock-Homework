package config

type Config struct{
	Server		ServerConfig		`mapstructure:"server"`
	ActivityRPC	ClientConfig		`mapstructure:"activity_rpc"`
	OrderRPC	ClientConfig		`mapstructure:"order_rpc"`
	Redis		RedisConfig		`mapstructure:"redis"`
}

// 这里为Hertz服务器的配置
type ServerConfig struct{
	Host 		string 	`mapstructure:"host"`
	Port 		int 	`mapstructure:"port"`
	LogLevel 	string 	`mapstructure:"log_level"`
}

// 这里是kitex client的配置
type ClientConfig struct{
	ServiceName string `mapstructure:"service_name"`
	TargetHost  string `mapstructure:"target_host"`
	TargetPort  int    `mapstructure:"target_port"`
	Timeout     int    `mapstructure:"timeout"`
}

// Redis 用于限流
type RedisConfig struct{
	Host 		string `mapstructure:"host"`
	Port 		int    `mapstructure:"port"`
	Password 	string `mapstructure:"password"`
	DB       	int    `mapstructure:"db"`
}
