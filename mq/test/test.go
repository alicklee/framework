package main

import (
	"strconv"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"

	"github.com/CloudcadeSF/Framework/mq"
)

func main() {
	mq.Start("rabbitmq")
	producer := mq.GlobalRMQ.Producers[1]
	s := "qqqq"
	go func() {
		for j := 0; j < 100; j++ {
			msg := s + strconv.Itoa(j)
			producer.SendMsg([]byte(msg))
		}
	}()
	log.Info("finished")
	select {}
}
