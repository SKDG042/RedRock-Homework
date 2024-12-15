package service

import (
	"Redrock/message-board/dao"
	"Redrock/message-board/model"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"time"
)

func Register(c context.Context, ctx *app.RequestContext) {
	var user model.User

	if err := ctx.BindAndValidate(&user); err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := dao.CreateUser(user); err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
		return
	}
	ctx.JSON(consts.StatusOK, utils.H{"message": "成功注册用户"})
}

func Login(c context.Context, ctx *app.RequestContext) {
	var LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BindAndValidate(&LoginRequest); err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	user, err := dao.GetUser(LoginRequest.Username)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
		return
	}

	if user.Password != LoginRequest.Password {
		ctx.JSON(consts.StatusUnauthorized, utils.H{"error": "账号或密码错误"})
		return
	} else {
		ctx.JSON(consts.StatusOK, utils.H{"message": "登录成功"})
	}
}
