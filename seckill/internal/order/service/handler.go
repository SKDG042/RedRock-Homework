package service

import (
	order "Redrock/seckill/kitex_gen/order"
	"context"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

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
