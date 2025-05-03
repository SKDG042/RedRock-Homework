package client

import(
	"time"
	"fmt"

	"github.com/cloudwego/kitex/client"

	"Redrock/seckill/internal/api/config"
	"Redrock/seckill/kitex_gen/activity/activityservice"
	"Redrock/seckill/kitex_gen/order/orderservice"
	"Redrock/seckill/kitex_gen/user/userservice"
)

type RPCClients struct{
	ActivityClient 	activityservice.Client
	OrderClient 	orderservice.Client
	UserClient 		userservice.Client
}

func NewRPCClients(cfg *config.Config) (*RPCClients, error){
	// 创建活动客户端
	activityClient, err := activityservice.NewClient(
		cfg.ActivityRPC.ServiceName,
		client.WithHostPorts(fmt.Sprintf("%s:%d", cfg.ActivityRPC.TargetHost, cfg.ActivityRPC.TargetPort)),
		client.WithRPCTimeout(time.Duration(cfg.ActivityRPC.Timeout)*time.Second),
	)

	if err != nil{
		return nil, fmt.Errorf("创建活动客户端失败：%v", err)
	}

	// 创建订单客户端
	orderClient, err := orderservice.NewClient(
		cfg.OrderRPC.ServiceName,
		client.WithHostPorts(fmt.Sprintf("%s:%d", cfg.OrderRPC.TargetHost, cfg.OrderRPC.TargetPort)),
		client.WithRPCTimeout(time.Duration(cfg.OrderRPC.Timeout)*time.Second),
	)

	if err != nil{
		return nil, fmt.Errorf("创建订单客户端失败：%v", err)
	}

	// 创建用户客户端
	userClient, err := userservice.NewClient(
		cfg.UserRPC.ServiceName,
		client.WithHostPorts(fmt.Sprintf("%s:%d", cfg.UserRPC.TargetHost, cfg.UserRPC.TargetPort)),
		client.WithRPCTimeout(time.Duration(cfg.UserRPC.Timeout)*time.Second),
	)

	if err != nil{
		return nil, fmt.Errorf("创建用户客户端失败：%v", err)
	}
	
	return &RPCClients{
		ActivityClient: activityClient,
		OrderClient: orderClient,
		UserClient: userClient,
	}, nil
}
