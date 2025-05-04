package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/app"

	"Redrock/seckill/internal/api/client"
	"Redrock/seckill/kitex_gen/order"
)

// OrderHandler 订单相关处理器
type OrderHandler struct{
	orderClients *client.RPCClients
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderClient *client.RPCClients) *OrderHandler{
	return &OrderHandler{
		orderClients: orderClient,
	}
}

// CreateOrder 创建秒杀订单
func (h *OrderHandler) CreateOrder(ctx context.Context, c *app.RequestContext){
	var req order.CreateOrderRequest
	if err := c.BindJSON(&req); err != nil{
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "请求的参数有误: " + err.Error(),
		})
		return
	}

	resp, err := h.orderClients.OrderClient.CreateOrder(ctx, &req)
	if err != nil{
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetOrder 获取秒杀订单详情
func (h *OrderHandler) GetOrder(ctx context.Context, c *app.RequestContext){
	userIDStr := c.Param("user_id")
	orderSn := c.Param("order_sn")

	if userIDStr == "" || orderSn == ""{
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "用户ID和订单号不能为空",
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil{
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "用户ID参数有误" + err.Error(),
		})
		return
	}

	req := &order.GetOrderRequest{
		UserID: userID,
		OrderSn: orderSn,
	}

	// 根据用户ID和订单号获取订单详情
	resp, err := h.orderClients.OrderClient.GetOrder(ctx, req)
	if err != nil{
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ListUserOrders 获取用户订单列表
func (h *OrderHandler) ListUserOrders(ctx context.Context, c *app.RequestContext){
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil{
		c.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "用户ID参数有误" + err.Error(),
		})
		return
	}

	var status order.OrderStatus = -1 // 默认获取所有状态的订单
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
		status = order.OrderStatus(statusInt)
	}

	req := &order.ListOrdersRequest{
		UserID: userID,
		Status: status,
	}

	// 根据用户ID和订单状态获取订单列表
	resp, err := h.orderClients.OrderClient.ListOrders(ctx, req)
	if err != nil{
		c.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "服务器内部错误: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}
