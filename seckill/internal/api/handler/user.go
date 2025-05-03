package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/app"

	"Redrock/seckill/internal/api/client"
	"Redrock/seckill/kitex_gen/user"
)


// UserHandler 用户相关处理器
type UserHandler struct {
	userClient *client.RPCClients
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userClient *client.RPCClients) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// Register 用户注册
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req user.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "请求的参数有误: " + err.Error(),
		})
		return
	}

	resp, err := h.userClient.UserClient.Register(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Login 用户登录
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req user.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "请求的参数有误: " + err.Error(),
		})
		return
	}

	resp, err := h.userClient.UserClient.Login(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}
