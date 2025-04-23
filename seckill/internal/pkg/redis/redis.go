package redis

import (
	"fmt"
	"log"
	"context"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

// 此处为共享的初始化Redis函数
func InitRedis(config *RedisConfig) error{
	var err error
	Client = redis.NewClient(&redis.Options{
		Addr 		: fmt.Sprintf("%s:%d",config.Host,config.Port),
		Password 	: config.Password,
		DB 			: config.DB,
		PoolSize 	: config.PoolSize,
	})
	
	// 接着连接并测试连通性
	ctx := context.Background()
	_,err = Client.Ping(ctx).Result()
	if err != nil{
		return fmt.Errorf("Redis连接失败：%w",err)
	}

	log.Printf("Redis连接成功")
	return nil
}

// GetRedis 用于返回Redis client的实例
func GetRedis() *redis.Client{
	if Client == nil{
		panic("您还没有初始化Redis client")
	}
	return Client
}

// CloseRedis 用于关闭 Redis 来凝结
func CloseRedis(){
	if Client != nil{
		if err := Client.Close(); err != nil{
			log.Printf("关闭Redis连接失败：%v", err)
		}else{
			log.Println("成功关闭Redis连接")
		}
	}
}