package api

import (
	"Redrock/message-board/service"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter() *server.Hertz {
	h := server.Default()

	h.POST("/register", service.Register)
	h.POST("/login", service.Login)
	h.POST("/message", service.Message)
	h.GET("/message", service.GetAllMessage)
	h.DELETE("/message", service.DeleteMessage)
	h.PUT("/user", service.UpdateUser)
	h.POST("/like", service.AddLike)
	h.DELETE("/like", service.DeleteLike)
	h.GET("/like", service.GetLike)

	return h
}
