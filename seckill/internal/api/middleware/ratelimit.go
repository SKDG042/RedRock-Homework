package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/redis/go-redis/v9"
)

// RatelimiterConfig 限流器配置
type RateLimiterConfig struct{
	Limit		int				// 每个周期允许的请求数
	Period		time.Duration	// 时间周期
	keyFunc		func(ctx *app.RequestContext) string	// 生成限流键的函数
}

// RedisRateLimiter 分布式限流中间件
func RedisRateLimiter(redisClient *redis.Client, config *RateLimiterConfig) app.HandlerFunc{
	return func(c context.Context, ctx *app.RequestContext){
		key := fmt.Sprintf("ratelimit:%s", config.keyFunc(ctx))
		
		// 1. 获取当前计数
		countKey := fmt.Sprintf("%s:count", key)
		expireKey := fmt.Sprintf("%s:expire", key)
		
		// 检查是否存在过期时间
		exists, err := redisClient.Exists(c, expireKey).Result()
		if err != nil {
			// 如果发生错误不限流,避免影响秒杀
			ctx.Next(c)
			return
		}
		
		// 不存在则创建
		if exists == 0 {
			redisClient.Set(c, countKey, 0, config.Period)
			redisClient.Set(c, expireKey, time.Now().Add(config.Period).Unix(), config.Period)
		}
		
		// 2. 获取计数并增加
		count, err := redisClient.Incr(c, countKey).Result()
		if err != nil {
			ctx.Next(c)
			return
		}
		
		// 3. 获取限制重置时间
		expireAt, err := redisClient.Get(c, expireKey).Int64()
		if err != nil {
			ctx.Next(c)
			return
		}
		resetAfter := time.Until(time.Unix(expireAt, 0))
		
		// 设置RateLimit相关的HTTP头
		// 限流上限 剩余可用请求数 限流重置时间
		ctx.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(int(max(0, int64(config.Limit)-count+1))))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(int64(resetAfter.Seconds()), 10))
		
		// 如果超出限制
		if count > int64(config.Limit) {
			ctx.JSON(consts.StatusTooManyRequests, map[string]any{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"wait":    resetAfter.Seconds(),
			})
			ctx.Abort()
			return
		}
		
		ctx.Next(c)
	}
}

// SeckillLimiter 秒杀接口限流器 (每秒100个请求)
func SeckillLimiter(redisClient *redis.Client) app.HandlerFunc{
	return RedisRateLimiter(redisClient, &RateLimiterConfig{
		Limit: 100,
		Period: time.Second,
		keyFunc: func(ctx *app.RequestContext) string{
			// 优先从请求头获取用户ID
			userID := string(ctx.GetHeader("x-User-ID"))
			
			// 如果没有用户ID，则使用客户端IP
			if userID == "" {
				userID = ctx.ClientIP()
			}
			
			return userID
		},
	})
}
