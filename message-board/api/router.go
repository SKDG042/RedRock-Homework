package api

import (
	"Redrock/message-board/service"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func InitRouter() *server.Hertz {
	h := server.Default()

	h.POST("/register", service.Register)

	return h
}
