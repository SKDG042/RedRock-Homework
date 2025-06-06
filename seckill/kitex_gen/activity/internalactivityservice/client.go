// Code generated by Kitex v0.13.1. DO NOT EDIT.

package internalactivityservice

import (
	activity "Redrock/seckill/kitex_gen/activity"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	DeductStock(ctx context.Context, req *activity.DeductStockRequest, callOptions ...callopt.Option) (r *activity.DeductStockResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kInternalActivityServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kInternalActivityServiceClient struct {
	*kClient
}

func (p *kInternalActivityServiceClient) DeductStock(ctx context.Context, req *activity.DeductStockRequest, callOptions ...callopt.Option) (r *activity.DeductStockResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeductStock(ctx, req)
}
