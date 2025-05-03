package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"Redrock/seckill/internal/order/data"
	"Redrock/seckill/internal/order/mq"
	"Redrock/seckill/internal/pkg/models"
	myRedis "Redrock/seckill/internal/pkg/redis"
	"Redrock/seckill/kitex_gen/activity"
	activityClient "Redrock/seckill/kitex_gen/activity/activityservice"
	internalClient "Redrock/seckill/kitex_gen/activity/internalactivityservice"
	order "Redrock/seckill/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{
	orderData 		*data.OrderData
	orderProducer 	*mq.OrderProducer
	activityClient 	activityClient.Client
	redisClient 	*redis.Client
	internalClient internalClient.Client
}

// NewOrderServiceImpl 创建服务实现实例
func NewOrderServiceImpl(producer *mq.OrderProducer, activityClient activityClient.Client) *OrderServiceImpl{
	internalActivityClient, err := internalClient.NewClient("Activity")
	if err != nil{
		panic(fmt.Sprintf("创建内部活动客户端失败：%v",err))
	}

	return &OrderServiceImpl{
		orderData: 		data.NewOrderData(),
		orderProducer: 	producer,
		activityClient: activityClient,
		redisClient: 	myRedis.GetRedis(),
		internalClient: internalActivityClient,
	}
}

// GenerateOrderSn 生成订单号
// 这里简单地用时间戳来表示
func generateOrderSn() string{
	
	return fmt.Sprintf("%d", time.Now().Unix())
}

// CreateOrder 创建订单
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (resp *order.CreateOrderResponse, err error) {
	response := &order.CreateOrderResponse{
		BaseResponse : &order.BaseResponse{},
	}

	if req.UserID <= 0 || req.ActivityID <= 0{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "输入参数错误"

		return response, nil
	}

	userID 		:= uint(req.UserID)
	activityID  := uint(req.ActivityID)

	// 生成订单号
	orderSn := generateOrderSn()

	// 1. 扣除库存
	deductResquest := &activity.DeductStockRequest{
		ActivityID:		req.ActivityID,
		UserID:			req.UserID,
		Count:			1, // 这里简单地扣1个
	}

	deductResponse, err := s.internalClient.DeductStock(ctx, deductResquest)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "扣除库存失败" + err.Error()

		return response, nil
	}

	// 判断扣除库存是否成功, 如果不成功则将扣除失败的响应返回给调用方
	if deductResponse.BaseResponse.Code != 0 || !deductResponse.Success{
		response.BaseResponse.Code = deductResponse.BaseResponse.Code
		response.BaseResponse.Msg  = deductResponse.BaseResponse.Msg
		
		return response, nil
	}

	// 2. 获取活动详情， 计算amount等
	activityRequest := &activity.GetActivityRequest{
		ActivityID:			int64(activityID),
	}

	activityResponse, err := s.activityClient.GetActivity(ctx, activityRequest)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "获取活动信息失败：" + err.Error()

		return response, nil
	}

	// 如果调用成功但未成功获取，如果活动已结束等
	if activityResponse.BaseResponse.Code != 0{
		response.BaseResponse.Code = activityResponse.BaseResponse.Code
		response.BaseResponse.Msg  = activityResponse.BaseResponse.Msg

		return response, nil
	}

	// 3.创建订单并写入数据库
	localOrder := &models.Order{
		OrderSn:			orderSn,
		UserID:				userID,
		ActivityID:			activityID,
		ProductID:			uint(activityResponse.Activity.ProductId),
		Amount:				activityResponse.Activity.SeckillPrice,   // 因为我们这里默认秒杀一件商品，所以seckillprice == amount
		Status:				models.StatusPending,
		CreateTime: 		time.Now(),	
	}

	err = s.orderData.Create(ctx, localOrder)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "创建订单失败：" + err.Error()

		return response, nil
	}

	// 4. 发送消息到mq，异步创建订单
	msg := &mq.OrderMessage{
		OrderSn:			orderSn,
		UserID:				userID,
		ActivityID:			activityID,
		ProductID:			uint(activityResponse.Activity.ProductId),
		Amount:				activityResponse.Activity.SeckillPrice,
	}
		
	err = s.orderProducer.Produce(msg)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "发送订单消息失败：" + err.Error()

		return response, nil
	}

	// 5. 构建返回的订单信息
	orderInfo := &order.OrderInfo{
		Id:					int64(localOrder.ID),
		OrderSn:			orderSn,
		UserID:				int64(userID),
		ActivityID:			int64(activityID),
		ProductID: 			int64(localOrder.ProductID),
		Amount:				localOrder.Amount,
		Status:				models.StatusPending,
		CreateTime: 		localOrder.CreatedAt.Unix(),	
	}

	response.OrderInfo = orderInfo
	response.BaseResponse.Code = 0
	response.BaseResponse.Msg  = "下单成功"
	
	return response, nil
}

// GetOrder 获取订单信息
func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderRequest) (resp *order.GetOrderResponse, err error) {
	response := &order.GetOrderResponse{
		BaseResponse: &order.BaseResponse{},
	}

	if req.UserID <= 0 || req.OrderSn == ""{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "输入参数错误"
	
		return response, nil
	}

	// 查询订单
	localOrder, err := s.orderData.GetByUserIDAndOrderSn(ctx, uint(req.UserID), req.OrderSn)
	if err != nil{
		response.BaseResponse.Code = 404
		response.BaseResponse.Msg  = "查询订单信息失败：" + err.Error()

		return response, nil
	}

	// 将orderStatus int转化为	api响应中的 enum
	var orderStatus order.OrderStatus

	switch localOrder.Status {
	case models.StatusPending:
		orderStatus = order.OrderStatus_PENDING
	case models.StatusCreated:
		orderStatus = order.OrderStatus_CREATED
	case models.StatusPaid:
		orderStatus = order.OrderStatus_PAID
	case models.StatusFailed:
		orderStatus = order.OrderStatus_FAILED
	case models.StatusCancelled:
		orderStatus = order.OrderStatus_CANCELLED
	default:
		orderStatus = order.OrderStatus_PENDING
	}

	orderInfo := &order.OrderInfo{
		Id:				int64(localOrder.ID),
		OrderSn:		localOrder.OrderSn,
		UserID:			int64(localOrder.UserID),
		ActivityID:		int64(localOrder.ActivityID),
		ProductID:		int64(localOrder.ProductID),
		Amount:			localOrder.Amount,
		Status:			orderStatus,
		CreateTime:		localOrder.CreatedAt.Unix(),
	}

	// 如果预加载成功
	if localOrder.Product.ID >0{
		orderInfo.ProductName = localOrder.Product.Name
	}

	response.OrderInfo = orderInfo
	response.BaseResponse.Code = 0
	response.BaseResponse.Msg  = "查询订单信息成功"
	
	return response, nil
}

// ListOrders 获取用户订单列表
func (s *OrderServiceImpl) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (resp *order.ListOrdersResponse, err error) {
	response := &order.ListOrdersResponse{
		BaseResponse: &order.BaseResponse{},
		Orders:	[]*order.OrderInfo{},
	}

	if req.UserID <= 0{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "用户ID不能为空"

		return response, nil
	}

	// 查询订单列表
	status := -1 // 默认为-1, 查询所有订单
	if req.Status != order.OrderStatus(-1){
		status = int(req.Status)
	}

	orders, total, err := s.orderData.ListByUserID(ctx, uint(req.UserID), status)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "查询订单列表失败：" + err.Error()

		return response, nil
	}

	// 构建返回的订单列表
	for _, o := range orders{
		// 转化orderStatus
		var orderStatus order.OrderStatus

		switch o.Status {
		case models.StatusPending:
			orderStatus = order.OrderStatus_PENDING
		case models.StatusCreated:
			orderStatus = order.OrderStatus_CREATED
		case models.StatusPaid:
			orderStatus = order.OrderStatus_PAID
		case models.StatusFailed:
			orderStatus = order.OrderStatus_FAILED
		case models.StatusCancelled:
			orderStatus = order.OrderStatus_CANCELLED
		default:
			orderStatus = order.OrderStatus_PENDING
		}

		orderInfo := &order.OrderInfo{
			Id: 			int64(o.ID),
			OrderSn: 		o.OrderSn,
			UserID: 		int64(o.UserID),
			ActivityID: 	int64(o.ActivityID),
			ProductID: 		int64(o.ProductID),
			Amount: 		o.Amount,
			Status: 		orderStatus,
			CreateTime: 	o.CreatedAt.Unix(),
		}

		if o.Product.ID > 0{
			orderInfo.ProductName = o.Product.Name
		}

		response.Orders = append(response.Orders, orderInfo)
	}
	
	response.Total = total
	response.BaseResponse.Code = 0
	response.BaseResponse.Msg  = "查询用户列表订单成功"

	return response, nil
}
