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
