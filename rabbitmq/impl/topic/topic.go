package topic

import (
	"encoding/json"

	"github.com/CloudcadeSF/Framework/iface/request"
	"github.com/CloudcadeSF/Framework/rabbitmq/impl"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/streadway/amqp"
)

/**
Topic的结构体
*/
type Topic struct {
	Rmq          *impl.RabbitMQ
	ExitChain    chan bool
	ExchangeType string
	Handler      func(msg []byte)
	Sender       func(req request.IRequest)
}

/**
创建一个新的Topic类型的MQ实体
*/
func NewTopic(queueName string, exchangeName string, routingKey string) *Topic {
	rmq := impl.NewRabbitMQ("rabbitmq")
	rmq.QueueName = queueName
	rmq.ExchangeName = exchangeName
	rmq.RoutingKey = routingKey
	t := &Topic{
		Rmq:          rmq,
		ExchangeType: "topic",
		ExitChain:    make(chan bool, 1),
	}
	return t
}

func (t *Topic) AddHandler(handler func(msg []byte)) {
	t.Handler = handler
}

func (t *Topic) AddSender(sender func(req interface{})) {
	panic("implement me")
}

/**
Topic发送消息方法
*/
func (t *Topic) SendMsg(msg []byte, routingKey string) error {
	r := t.Rmq
	ch, err := r.Conn.Channel()
	if err != nil {
		log.Panic("Get channel error in topic", err)
	}
	if err := ch.Publish(r.ExchangeName, routingKey, false, false, amqp.Publishing{
		UserId: "",
		AppId:  "",
		Body:   msg,
	}); err != nil {
		return err
	}
	log.Infoln("rabbitmq has sent message")
	return nil
}

/**
Topic发送fanout消息方法
*/
func (t *Topic) SendListMsg(msg []byte, pids []string) error {
	pidStrings, err2 := json.Marshal(pids)
	if err2 != nil {
		return err2
	}
	r := t.Rmq
	ch, err := r.Conn.Channel()
	if err != nil {
		log.Panic("Get channel error in topic", err)
	}
	if err := ch.Publish(r.ExchangeName, r.RoutingKey, false, false, amqp.Publishing{
		UserId: string(pidStrings),
		Body:   msg,
	}); err != nil {
		return err
	}
	log.Infoln("rabbitmq has sent message")
	return nil
}

/**
Topic发送fanout消息方法
*/
func (t *Topic) FanoutMsg(msg []byte) error {
	r := t.Rmq
	ch, err := r.Conn.Channel()
	if err != nil {
		log.Panic("Get channel error in topic", err)
	}
	if err := ch.Publish(r.ExchangeName, r.RoutingKey, false, false, amqp.Publishing{
		Body: msg,
	}); err != nil {
		return err
	}
	log.Infoln("rabbitmq has sent message")
	return nil
}

/**
Topic消费消息的方法
*/
func (t *Topic) RecvMsg() error {
	//消费消息
	ch, err := t.Rmq.Conn.Channel()
	if err != nil {
		return err
	}
	if msgs, err := ch.Consume(t.Rmq.QueueName, "", false, false, false, false, nil); err != nil {
		return err
	} else {
		forever := make(chan bool)
		for data := range msgs {
			//TODO 消费消息的逻辑部分
			log.Infoln(string(data.Body))
			t.Handler(data.Body)
			ch.Ack(data.DeliveryTag, true)
		}
		<-forever
	}
	return nil
}

/**
Topic消费消息并转发
*/
func (t *Topic) RecvAndRelay() error {
	panic("implement me")
}

/**
启动Topic的消费服务
*/
func (t *Topic) Start() error {
	if err := t.Rmq.ExchangeDeclare(t.Rmq.ExchangeName, t.ExchangeType); err != nil {
		return err
	}
	if err := t.Rmq.QueueDeclare(t.Rmq.QueueName); err != nil {
		return err
	}
	if err := t.Rmq.BindQueue(t.Rmq.QueueName, t.Rmq.ExchangeName, t.Rmq.RoutingKey); err != nil {
		return err
	}
	return nil
}

/**
停止监听
*/
func (t *Topic) Stop() {
	t.ExitChain <- true
	//close(t.ExitChain)
}
