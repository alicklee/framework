package iface

import "github.com/streadway/amqp"

/**
RabbitMQ的接口
*/
type IRabbitMQ interface {
	//启动一个RabbitMQ实例
	Start() error
	//停止RabbitMQ
	Stop() error
	//申明exchange
	ExchangeDeclare(exchangeName string, exchangeType string) error
	//声明一个延时交换机
	DelayExchangeDeclare(exchangeName string, exchangeType string,arg amqp.Table) error
	//申明队列
	QueueDeclare(queueName string) error
	//绑定queue
	BindQueue(queueName string, exchangeName string, routingKey string) error
}
