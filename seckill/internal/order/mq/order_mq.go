package mq

import (
	"context"
	"fmt"
	"encoding/json"
	"log"
	"time"

	"Redrock/seckill/internal/pkg/mq"
	"Redrock/seckill/internal/pkg/models"
	"Redrock/seckill/internal/order/data"
)

// OrderMessage 订单消息
type OrderMessage struct{
	OrderSn		string		`json:"order_sn"`
	UserID		uint		`json:"user_id"`
	ActivityID	uint		`json:"activity_id"`
	ProductID	uint		`json:"product_id"`
	Amount		float64		`json:"amount"`
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
		return nil, fmt.Errorf("创建订单消费者失败：%w", err)
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

// handlerOrderMessage 处理订单消息
func (c *OrderConsumer) handlerOrderMessage(body []byte) error{
	var msg OrderMessage

	err := json.Unmarshal(body, &msg)
	if err != nil{
		return fmt.Errorf("反序列化消息失败：%w", err)
	}

	log.Printf("收到订单消息：%v", msg)

	// 创建订单
	order := &models.Order{
		OrderSn:		msg.OrderSn,
		UserID:		msg.UserID,
		ActivityID:	msg.ActivityID,
		ProductID:	msg.ProductID,
		Amount:		msg.Amount,
		Status:		models.StatusCreated,
		CreateTime:	time.Now(),
	}

	// 写入数据库
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = c.orderData.Create(ctx, order)
	if err != nil{
		log.Printf("创建订单失败：%v", err)
		// 更新订单状态为失败
		errUpdate := c.orderData.UpdateStatus(ctx, msg.OrderSn, models.StatusFailed)
		if errUpdate != nil{
			log.Printf("更新订单状态失败：%v", errUpdate)
		}
		return fmt.Errorf("创建订单失败：%w", err)
	}

	log.Printf("订单创建成功：%v", order)

	return nil
}

// StartConsume 开始消费订单消息
func (c *OrderConsumer) StartConsume() error{
	err := c.rabbitmq.ConsumeMessage(c.handlerOrderMessage)

	return err
}
