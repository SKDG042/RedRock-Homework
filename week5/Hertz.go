package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type PingRequest struct {
	Message string `json:"message"`
}

func main() {
	h := server.Default() // 创建engine
	h.POST("/ping", func(c context.Context, ctx *app.RequestContext) {
		var req PingRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		ctx.JSON(consts.StatusOK, utils.H{"message": req.Message})
	}) // 定义一个/ping路由，支持POST方法
	h.Spin()
}
