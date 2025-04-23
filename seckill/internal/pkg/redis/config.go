package redis

type RedisConfig struct{
	Host 		string 	`mapstructure:"host"`
	Port 		int 	`mapstructure:"port"`
	Password 	string 	`mapstructure:"password"`
	DB 			int 	`mapstructure:"db"`
	PoolSize 	int 	`mapstructure:"pool_size"` // 连接池大小，为了提高并发性能 默认100
}