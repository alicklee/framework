package channel

import (
	"github.com/CloudcadeSF/Framework/mq/exchange"
	"github.com/CloudcadeSF/Framework/mq/queue"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/streadway/amqp"
)

/**
RabbitMQ的channel的结构体
*/
type Channel struct {
	//通道名称
	Name string
	//通道的交换机
	Exchange *exchange.Exchange
	//消息队列
	Queue *queue.Queue
	//MQ的channel句柄
	Ch *amqp.Channel
	//MQ的连接句柄
	Conn *amqp.Connection
}

/**
创建一个新的channel实体
*/
func NewChannel(name string, ex *exchange.Exchange, q *queue.Queue, conn *amqp.Connection, ch *amqp.Channel) *Channel {
	channel := &Channel{
		Name:     name,
		Exchange: ex,
		Queue:    q,
		Conn:     conn,
		Ch:       ch,
	}
	return channel
}

/**
声明一个交换机
*/
func (c *Channel) DeclareExchange() error {
	err := c.Ch.ExchangeDeclare(
		c.Exchange.Name,
		c.Exchange.Type,
		c.Exchange.Durable,
		c.Exchange.AutoDeleted,
		c.Exchange.Internal,
		c.Exchange.NoWait,
		nil)
	if err != nil {
		log.Errorln("When channel declare exchange has a error", err)
		return err
	}
	return nil
}

/**
通过一个routing key 绑定queue
*/
func (c *Channel) BindQueue() error {
	if _, err := c.Ch.QueueDeclare(c.Queue.Name, c.Queue.Durable, c.Queue.AutoDel, c.Queue.Exclusive, c.Queue.NoWait, nil); err != nil {
		log.Errorln("When channel create a queue by routing key error", err)
	}
	if err := c.Ch.QueueBind(c.Queue.Name, c.Queue.RoutingKey, c.Exchange.Name, false, nil); err != nil {
		log.Errorln("When channel bing a queue by routing key error", err)
		return err
	}
	return nil
}
