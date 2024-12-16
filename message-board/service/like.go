package service

import (
	"Redrock/message-board/dao"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

func AddLike(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.PostForm("user_id")
	messageIDStr := ctx.PostForm("message_id")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"error": "用户ID不存在"})
		return
	}

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "留言ID不存在"})
		return
	}

	if err := dao.AddLike(userID, messageID); err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "点赞失败"})
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"message": "成功点赞"})
}

func DeleteLike(c context.Context, ctx *app.RequestContext) {
	idStr := ctx.PostForm("user_id")
	messageIDStr := ctx.PostForm("message_id")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{"error": "用户ID不存在"})
		return
	}

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "留言ID不存在"})
		return
	}

	if err := dao.DeleteLike(userID, messageID); err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "取消点赞失败"})
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"message": "成功取消点赞"})
}

func GetLike(c context.Context, ctx *app.RequestContext) {
	messageIDStr := ctx.Query("message_id")

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "留言ID不存在"})
		return
	}

	likes, err := dao.GetLike(messageID)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "获取点赞数失败"})
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"likes": likes})
}
