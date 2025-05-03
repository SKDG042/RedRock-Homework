package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"Redrock/seckill/internal/order/data"
	"Redrock/seckill/internal/pkg/models"
	"Redrock/seckill/internal/pkg/mq"
)

// OrderMessage 订单消息
type OrderMessage struct{
	OrderSn		string		`json:"order_sn"`
	UserID		uint		`json:"user_id"`
	ActivityID	uint		`json:"activity_id"`
	ProductID	uint		`json:"product_id"`
	Amount		float64		`json:"amount"`
	Price       float64     `json:"price"`
	Quantity    int         `json:"quantity"`
}

// OrderProducer 订单消息生产者
type OrderProducer struct{
	rabbitmq *mq.RabbitMQ
}

// OrderConsumer 订单消息消费者
type OrderConsumer struct{
	rabbitmq *mq.RabbitMQ
	orderData *data.OrderData
}

// NewOrderProducer 创建订单消息生产者
func NewOrderProducer(config *mq.RabbitMQConfig) (*OrderProducer, error){
	localRabbitMQ, err := mq.NewRabbitMQ(config)
	if err != nil{
		return nil, fmt.Errorf("创建订单生产者失败：%w", err)
	}

	return &OrderProducer{
		rabbitmq: localRabbitMQ,
	}, nil
}

// Close 关闭连接
func (p *OrderProducer) Close(){
	if p.rabbitmq != nil{
		p.rabbitmq.Close()
	}
}

// Produce 生产订单消息
func (p *OrderProducer) Produce(message *OrderMessage) error{

	if message.OrderSn == "" {
        return fmt.Errorf("订单号不能为空")
    }
    
    if message.UserID == 0 {
        return fmt.Errorf("用户ID不能为0")
    }
    
    if message.ActivityID == 0 {
        return fmt.Errorf("活动ID不能为0")
    }
    
    if message.ProductID == 0 {
        return fmt.Errorf("商品ID不能为0")
    }

	data, err := json.Marshal(message)
	if err != nil{
		return fmt.Errorf("序列化消息失败：%w", err)
	}

	err = p.rabbitmq.PublishMessage(data)

	return err
}

// NewOrderConsumer 创建订单消息消费者
func NewOrderConsumer(config *mq.RabbitMQConfig, orderData *data.OrderData) (*OrderConsumer, error){
	rabbit, err := mq.NewRabbitMQ(config)
	if err != nil{
		return nil, fmt.Errorf("创建订单消息消费者失败：%w", err)
	}

	return &OrderConsumer{
		rabbitmq: 	rabbit,
		orderData: 	orderData,
	}, nil
}

// Close 关闭连接
func (c *OrderConsumer) Close(){
	if c.rabbitmq != nil{
		c.rabbitmq.Close()
	}
}

// handlerOrderMessage 处理订单消息(仅负责确认收到消息更新状态)
func (c *OrderConsumer) handlerOrderMessage(body []byte) error{
	var msg OrderMessage

	err := json.Unmarshal(body, &msg)
	if err != nil{
		return fmt.Errorf("反序列化消息失败：%w", err)
	}

	log.Printf("收到订单消息：%v", msg)

	// 检查订单是否存在
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exist, err := c.orderData.GetByOrderSn(ctx, msg.OrderSn)
	if err != nil{
		return fmt.Errorf("检查订单是否存在失败：%w", err)
	}

	if exist == nil{
		return fmt.Errorf("订单不存在")
	}

	if exist.Status == models.StatusPending {
		log.Printf("开始更新订单%v状态", msg.OrderSn)

		err = c.orderData.UpdateStatus(ctx, msg.OrderSn, models.StatusCreated)
		if err != nil {
			return fmt.Errorf("更新订单状态失败：%w", err)
		}
		
		log.Printf("订单%v状态更新成功", msg.OrderSn)
	} else {
		log.Printf("订单%v当前状态不是Pending(状态码:%d)，无需更新", msg.OrderSn, exist.Status)
	}

	return nil
}

// StartConsume 开始消费订单消息
func (c *OrderConsumer) StartConsume() error{
	err := c.rabbitmq.ConsumeMessage(c.handlerOrderMessage)

	return err
}
