package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strings"
)

func main() {
	h := server.Default() //创建服务器，默认端口为8888

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.GET("/echo", func(c context.Context, ctx *app.RequestContext) {
		message := ctx.Query("message")
		message = strings.ToLower(message)
		ctx.JSON(consts.StatusOK, utils.H{"message": message})
	})

	h.Spin() //启动服务器
}
