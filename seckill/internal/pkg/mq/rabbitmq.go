package mq

import(
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// 封装RabbitMQ
type RabbitMQ struct{
	connection 	*amqp.Connection
	channel		*amqp.Channel
	config		*RabbitMQConfig
}

// NewRabbitMQ 创建RabbitMQ实例
func NewRabbitMQ(config *RabbitMQConfig) (*RabbitMQ, error){
	// 连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Password, config.Host, config.Port))
	if err != nil{

		return nil, fmt.Errorf("连接RabbitMQ失败: %w", err)
	}

	// 创建通道
	ch, err := conn.Channel()
	if err != nil{
		conn.Close()
		
		return nil, fmt.Errorf("创建通道失败: %w", err)
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		config.ExchangeName,	// 交换机名称
		"direct",				// 交换机类型
		true,					// 是否持久化
		false,					// 是否自动删除
		false,					// 是否内部交换机
		false,					// 是否阻塞
		nil,					// 其他参数
	)
	if err != nil{
		ch.Close()
		conn.Close()

		return nil, fmt.Errorf("声明交换机失败: %w", err)
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		config.QueueName,	// 队列名称
		true,				// 是否持久化
		false,				// 是否自动删除
		false,				// 是否排他
		false,				// 是否阻塞
		nil,				// 其他参数
	)
	if err != nil{
		ch.Close()
		conn.Close()

		return nil, fmt.Errorf("声明队列失败: %w", err)
	}

	// 将队列绑定到交换机
	err = ch.QueueBind(
		config.QueueName,		// 队列名称
		config.RoutingKey,		// 路由键
		config.ExchangeName,	// 交换机名称
		false,					// 是否阻塞
		nil,					// 其他参数
	)
	if err != nil{
		ch.Close()
		conn.Close()

		return nil, fmt.Errorf("绑定队列到交换机失败: %w", err)
	}

	return &RabbitMQ{
		connection:	conn,
		channel:	ch,
		config:		config,
	}, nil
}

// Close 关闭连接
func (r *RabbitMQ) Close(){
	if r.channel != nil{
		r.channel.Close()
	}

	if r.connection != nil{
		r.connection.Close()
	}
}

// PublishMessage 发布消息
func (r *RabbitMQ) PublishMessage(body []byte) error{
	err := r.channel.Publish(
		r.config.ExchangeName,	// 交换机名称
		r.config.RoutingKey,	// 路由键名称
		false,					// 是否强制发送
		false,					// 是否立即发送
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久化消息
			ContentType:  "application/json",
			Body:		  body,
		},
	)

	return fmt.Errorf("发布消息失败：%w", err)
}

// ConsumeMessage 消费消息
func (r *RabbitMQ) ConsumeMessage(handler func([]byte) error) error{
	// Quality of Service 服务质量
	// 设置QoS为1，表示每次只处理一条消息
	err := r.channel.Qos(1, 0, false) // 预取数量，大小限制，是否全局
	if err != nil{
		return fmt.Errorf("设置QoS失败: %w", err)
	}

	// 注册消费者
	msgs, err := r.channel.Consume(
		r.config.QueueName,	// 队列名称
		"",					// 消费者名称(自动生成)
		false,				// 是否自动应答
		false,				// 是否排他
		false,				// 是否开启本地模型
		false,				// 是否阻塞
		nil,
	)

	if err != nil{
		return fmt.Errorf("注册消费者失败：%w", err)
	}

	go func(){
		for msg := range msgs{
			log.Printf("收到消息：%s", msg.Body)

			err := handler(msg.Body)
			if err != nil{
				log.Printf("处理消息失败：%v", err)
				
				// 如果处理消息失败，则拒绝消息重新入对，重新尝试处理消息
				msg.Nack(false,true)
			}else{
				// 处理消息成功, 确认消息
				msg.Ack(false)
			}
		}
	}()

	log.Printf("消费者成功启动，等待消息中...")

	return nil
}
