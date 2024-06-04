package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/CloudcadeSF/Framework/mq/channel"
	"github.com/CloudcadeSF/Framework/mq/message"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/streadway/amqp"
)

type Consumer struct {
	//生产者通道
	Name string
	//Channel
	Channel *channel.Channel
	//生产者编号
	Id int
}

func NewConsumer(id int, name string, ch *channel.Channel) *Consumer {
	c := &Consumer{
		Name:    name,
		Channel: ch,
		Id:      id,
	}
	return c
}

func (c *Consumer) ConsumerMsg() <-chan amqp.Delivery {
	msgs, err := c.Channel.Ch.Consume(
		c.Channel.Queue.Name, // queue
		c.Name,               // consumer
		true,                 // auto ack
		false,                // exclusive
		false,                // no local
		false,                // no wait
		nil,                  // args
	)
	if err != nil {
		log.Errorln("consumer msg error", err)
	}
	return msgs
}

func (c *Consumer) MsgHandler() {
	msgs := c.ConsumerMsg()
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			msg := message.Message{}
			json.Unmarshal(d.Body, &msg)
			fmt.Printf("The data is %s, The msgID is %d\n", msg.Body, msg.MsgId)
		}
	}()
	<-forever
}
