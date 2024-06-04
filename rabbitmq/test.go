package main

import (
	"fmt"
	"github.com/CloudcadeSF/Framework/rabbitmq/impl/message"
	"log"
	"math/rand"
	"time"
)

func handlerMsg(i interface{}) {
	fmt.Println(string(i.([]byte)))
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func main() {
	//go func() {
	//	delayMsg := delay.NewDelay("game4.delay.queue","game4.delay.server","game4.server.delay.1")
	//	if err := delayMsg.Start(); err != nil {
	//		log.Error(err)
	//		return
	//	}
	//	go func() {
	//		i := 0
	//		for{
	//			if i >= 1000{
	//				log.Infoln("msg send finished")
	//				break
	//			}
	//			rand.Seed(time.Now().Unix())
	//			delayTime := int64(random(10000, 99999))
	//			delayMsg.SendMsg([]byte("sssss222"),delayTime)
	//			i++
	//
	//		}
	//	}()
	//	delayMsg.RecvMsg(handlerMsg)
	//}()
	//select {}
	m := message.Message{
		Conn:   nil,
		Body:   nil,
		UserId: "sss",
	}
	b := m.Encode()
	n := message.Decode(b)
	log.Println(n.GetUserId())
}
