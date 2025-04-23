package mq

// message发给exchange,再根据routingKey 分拣到 对应的queue
type RabbitMQConfig struct{
	Host			string		`mapstructure:"host"`
	Port			int			`mapstructure:"port"`
	User			string		`mapstructure:"user"`
	Password		string		`mapstructure:"password"`
	ExchangeName	string		`mapstructure:"exchange_name"`
	QueueName		string		`mapstructure:"queue_name"`
	RoutingKey		string		`mapstructure:"routing_key"`
}

type MQConfig struct{
	Type	string	`mapstructure:"type"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}
