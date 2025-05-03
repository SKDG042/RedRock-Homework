package handler

import (
	"context"
	"strconv"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/app"

	"Redrock/seckill/internal/api/client"
	"Redrock/seckill/kitex_gen/activity"
)

// ActivityHandler 活动相关处理器
type ActivityHandler struct{
	activityClients *client.RPCClients
}

// NewActivityHandler 创建活动处理器
func NewActivityHandler(activityClient *client.RPCClients) *ActivityHandler{
	return &ActivityHandler{
		activityClients: activityClient,
	}
}

// CreateActivity 创建秒杀活动
func (h *ActivityHandler) CreateActivity(ctx context.Context, c *app.RequestContext){
	var req activity.CreateActivityRequest
	if err := c.BindJSON(&req); err != nil{
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "请求的参数有误: " + err.Error(),
		})
		return
	}

	resp, err := h.activityClients.ActivityClient.CreateActivity(ctx, &req)
	if err != nil{
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ListActivities 获取秒杀活动列表
func (h *ActivityHandler) ListActivities(ctx context.Context, c *app.RequestContext){
	status := int32(-1) // 默认获取所有状态的活动

	statusStr := c.Query("status")
	if statusStr != ""{
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil{
			c.JSON(consts.StatusBadRequest, map[string]any{
				"code":    400,
				"message": "status参数有误" + err.Error(),
			})
			return
		}
		status = int32(statusInt)
	}

	req := &activity.GetActivityListRequest{
		Status: status,
	}

	// 根据status获取活动列表
	resp, err := h.activityClients.ActivityClient.GetActivityList(ctx, req)
	if err != nil{
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetActivity 获取秒杀活动详情
func (h *ActivityHandler) GetActivity(ctx context.Context, c *app.RequestContext){
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil{
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "id参数有误" + err.Error(),
		})
		return
	}

	req := &activity.GetActivityRequest{
		ActivityID: id,
	}

	resp, err := h.activityClients.ActivityClient.GetActivity(ctx, req)
	if err != nil{
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}
