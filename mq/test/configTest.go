package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/CloudcadeSF/Framework/mq/channel"

	"github.com/CloudcadeSF/Framework/mq/exchange"

	"github.com/CloudcadeSF/shop-heroes-legends-common/ccutil/value"

	"github.com/CloudcadeSF/shop-heroes-legends-common/apollo"
	"github.com/zouyx/agollo/v3"
)

func main() {
	//运行时要设置两个环境变量，APP_ENV=dev;CFG_SERVER=http://118.24.2.131:18080
	if err := apollo.InitDefault("rabbitmq"); err != nil {
		log.Panicf(err.Error())
	}
	ip := agollo.GetValue("dev.host")
	port := agollo.GetIntValue("dev.port", 0)
	username := agollo.GetValue("dev.username")
	password := agollo.GetValue("dev.password")
	var list []*exchange.Exchange
	serverListJson := agollo.GetStringValue("dev.exchanges", value.EmptyString)
	json.Unmarshal([]byte(serverListJson), &list)

	var channelsList []*channel.Channel
	channelJson := agollo.GetStringValue("dev.channels", value.EmptyString)
	json.Unmarshal([]byte(channelJson), &channelsList)

	fmt.Printf("ip is %s, port is %d, usernamne is %s, password id %s\n", ip, port, username, password)
	fmt.Println(channelsList[0].Exchange.Name)
	fmt.Println(channelsList[0].Name)
	fmt.Println(channelsList[0].Queue.Name)
}
