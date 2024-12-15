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

	return h
}
