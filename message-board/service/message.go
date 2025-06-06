package service

import (
	"Redrock/message-board/dao"
	"Redrock/message-board/model"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"time"
)

func Message(_ context.Context, ctx *app.RequestContext) {
	var message model.Message

	if err := ctx.BindAndValidate(&message); err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	if err := dao.CreateMessage(message); err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"message": "成功发表留言"})

}

func GetAllMessage(_ context.Context, ctx *app.RequestContext) {
	messages, err := dao.GetAllMessages()
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
		return
	}
	ctx.JSON(consts.StatusOK, messages)
}

func DeleteMessage(_ context.Context, ctx *app.RequestContext) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
		return
	}

	err = dao.DeleteMessage(id)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"message": "成功删除留言"})
}
