package main

import (
	"encoding/json"
	"fmt"
	"github.com/CloudcadeSF/Framework/mq"
	"github.com/CloudcadeSF/Framework/mq/message"
)

func main() {
	mq.Start("rabbitmq", "gateway")
	c := mq.GlobalRMQ.Consumers[0]
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
