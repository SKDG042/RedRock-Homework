package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{ // 使用上下文c.JSON()方法完成JSON响应
			"message": "poooooooong",
		})
	}) // 定义一个/ping路由，支持GET方法
	r.Run(":8888") // 监听并在 0.0.0.0:8080 上启动服务
}
