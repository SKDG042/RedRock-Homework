package utils

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Middleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.JSON(consts.StatusUnauthorized, utils.H{"error": "未登录"})
			ctx.Abort()
			return
		}

		claims, err := ParseToken(auth)
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, utils.H{"error": "Token无效"})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)

		ctx.Next(c)
	}
}
