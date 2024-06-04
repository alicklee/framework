package delay

import (
	"github.com/CloudcadeSF/Framework/rabbitmq/impl"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/streadway/amqp"
)

type Delay struct {
	Rmq          *impl.RabbitMQ
	ExitChain    chan bool
	ExchangeType string
	header       string
}

func NewDelay(queueName string, exchangeName string, routingKey string) *Delay {
	rmq := impl.NewRabbitMQ("rabbitmq")
	rmq.QueueName = queueName
	rmq.ExchangeName = exchangeName
	rmq.RoutingKey = routingKey
	d := &Delay{
		Rmq:          rmq,
		ExchangeType: "topic",
		ExitChain:    make(chan bool, 1),
	}
	return d
}

/**
发送延时消息
@params
msg   消息内容
delay 延时的时间单位为毫秒
*/
func (d *Delay) SendMsg(msg []byte, delay int64) error {
	r := d.Rmq
	ch, err := r.Conn.Channel()
	if err != nil {
		log.Panic("Get channel error in topic", err)
	}
	headers := make(amqp.Table)
	if delay != 0 {
		headers["x-delay"] = delay
	}
	if err := ch.Publish(r.ExchangeName, r.RoutingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/pain",
		Body:         msg,
		Headers:      headers,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Delay) RecvMsg(msgHandler func(interface{})) error {
	//消费消息
	ch, err := d.Rmq.Conn.Channel()
	if err != nil {
		return err
	}
	if msgs, err := ch.Consume(d.Rmq.QueueName, "", false, false, false, false, nil); err != nil {
		return err
	} else {
		forever := make(chan bool)
		for data := range msgs {
			//TODO 消费消息的逻辑部分
			msgHandler(data.Body)
			ch.Ack(data.DeliveryTag, true)
		}
		<-forever
	}
	return nil
}

func (d *Delay) RecvAndRelay() error {
	panic("implement me")
}

func (d *Delay) Start() error {
	//声明一个延时队列的交换机
	if err := d.Rmq.DelayExchangeDeclare(d.Rmq.ExchangeName, d.ExchangeType); err != nil {
		return err
	}
	//声明一个延时队列的queue
	if err := d.Rmq.QueueDeclare(d.Rmq.QueueName); err != nil {
		return err
	}
	//绑定queue到对应的交换机上面
	if err := d.Rmq.BindQueue(d.Rmq.QueueName, d.Rmq.ExchangeName, d.Rmq.RoutingKey); err != nil {
		return err
	}
	return nil
}

func (d *Delay) Stop() {
	d.ExitChain <- true
}
