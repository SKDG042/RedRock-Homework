package api

import (
	"Redrock/message-board/service"
	"Redrock/message-board/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter() *server.Hertz {
	h := server.Default()

	h.POST("/register", service.Register)
	h.POST("/login", service.Login)
	h.GET("/like", service.GetLike)
	h.GET("/message", service.GetAllMessage)

	auth := h.Group("/")
	auth.Use(utils.Middleware())
	{
		auth.POST("/message", service.Message)
		auth.DELETE("/message", service.DeleteMessage)
		auth.PUT("/user", service.UpdateUser)
		auth.POST("/like", service.AddLike)
		auth.DELETE("/like", service.DeleteLike)
	}
	return h
}
