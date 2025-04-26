package service

import (
	activity "Redrock/seckill/kitex_gen/activity"
	order "Redrock/seckill/kitex_gen/order"
	"context"
)

// InternalActivityServiceImpl implements the last service interface defined in the IDL.
type InternalActivityServiceImpl struct{}

// DeductStock implements the InternalActivityServiceImpl interface.
func (s *InternalActivityServiceImpl) DeductStock(ctx context.Context, req *activity.DeductStockRequest) (resp *activity.DeductStockResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (resp *order.CreateOrderResponse, err error) {
	// TODO: Your code here...
	return
}

// GetOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderRequest) (resp *order.GetOrderResponse, err error) {
	// TODO: Your code here...
	return
}

// ListOrders implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (resp *order.ListOrdersResponse, err error) {
	// TODO: Your code here...
	return
}
