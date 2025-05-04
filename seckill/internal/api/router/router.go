package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"Redrock/seckill/internal/api/client"
	"Redrock/seckill/internal/api/handler"
	"Redrock/seckill/internal/api/middleware"
	"Redrock/seckill/internal/pkg/redis"

)

func SetupRouter(h *server.Hertz, clients *client.RPCClients){
	// 获取Redis实例
	redisClient := redis.GetRedis()
	
	// 创建处理器
	userHandler := handler.NewUserHandler(clients)
	activityHandler := handler.NewActivityHandler(clients)
	orderHandler := handler.NewOrderHandler(clients)

	// API路由
	api := h.Group("/api")

	// 用户相关路由
	userGroup := api.Group("/user")
	{
	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)
	}
	// 活动相关路由
	activityGroup := api.Group("/activity")
	{
		activityGroup.POST("/create", activityHandler.CreateActivity)
		activityGroup.GET("/list", activityHandler.ListActivities)
		activityGroup.GET("/detail/:id", activityHandler.GetActivity)
	}

	// 订单相关路由
	orderGroup := api.Group("/order")
	{
		orderGroup.POST("/seckill", middleware.SeckillLimiter(redisClient), orderHandler.CreateOrder)      		// 秒杀接口
		orderGroup.GET("/status", orderHandler.GetOrder)     		// 查询订单状态
		orderGroup.GET("/list/:user_id", orderHandler.ListUserOrders) 	// 获取用户订单列表
	}
}
