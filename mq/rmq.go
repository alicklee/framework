package mq

import (
	"encoding/json"
	"fmt"

	"github.com/CloudcadeSF/Framework/mq/channel"
	"github.com/CloudcadeSF/Framework/mq/message"

	"github.com/CloudcadeSF/shop-heroes-legends-common/ccutil/value"

	"github.com/CloudcadeSF/Framework/mq/connection"
	"github.com/CloudcadeSF/Framework/mq/consumer"
	"github.com/CloudcadeSF/Framework/mq/enum"
	"github.com/CloudcadeSF/Framework/mq/producer"
	"github.com/CloudcadeSF/shop-heroes-legends-common/apollo"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/zouyx/agollo/v3"
)

var GlobalRMQ = RMQ{}

type RMQ struct {
	Producers map[int]*producer.Producer
	Consumers map[int]*consumer.Consumer
}

/**
启动rabbitmq的服务
*/
func Start(configName string, servername string) {
	//分别初始化生产者和消费者
	GlobalRMQ.initProducers(configName)
	GlobalRMQ.initConsumers(configName)
	//启动所有的消费者监听
	RunConsumers(servername)
}

/**
初始化所有的生产者
*/
func (r *RMQ) initProducers(configName string) {
	//从配置中获取rabbitmq相关的配置参数
	host, port, username, password, channelList := r.getChannelConf(configName)
	//所有的生产者复用一个连接
	conn, err := connection.NewConnection(host, port, username, password).Start()
	if err != nil {
		log.Errorln("Connection create failed", err)
	}
	for _, chconf := range channelList {
		ch, _ := conn.Channel()
		//从配置中获取channel的参数，创建一个Channel的实体
		initChannel := channel.NewChannel(chconf.Name, chconf.Exchange, chconf.Queue, conn, ch)
		//创建一个生产者
		r.initOneProducer(initChannel)
	}
}

/**
初始化所有的消费者
*/
func (r *RMQ) initConsumers(configName string) {
	host, port, username, password, channelList := r.getChannelConf(configName)
	//所有的消费者复用一个连接
	conn, err := connection.NewConnection(host, port, username, password).Start()
	if err != nil {
		log.Errorln("Connection create failed", err)
	}
	i := 0
	for _, chconf := range channelList {
		ch, _ := conn.Channel()
		initChannel := channel.NewChannel(chconf.Name, chconf.Exchange, chconf.Queue, conn, ch)
		r.initOneConsumer(i, initChannel)
		i++
	}
}

/**
初始化单个生产者
*/
func (r *RMQ) initOneProducer(ch *channel.Channel) {
	r.Producers = make(map[int]*producer.Producer)
	p := producer.NewProducer(enum.GameServer1, ch.Exchange.Name, ch)
	r.Producers[p.Id] = p
}

/**
初始化单个消费者
*/
func (r *RMQ) initOneConsumer(id int, ch *channel.Channel) error {
	r.Consumers = make(map[int]*consumer.Consumer)

	//创建一个消费者
	c := consumer.NewConsumer(id, ch.Exchange.Name, ch)
	//声明一个exchange
	if err := ch.DeclareExchange(); err != nil {
		log.Errorln("Declare exchange error", err)
		return err
	}
	//绑定一个queue
	if err := ch.BindQueue(); err != nil {
		log.Errorln("BindQueue error", err)
		return err
	}
	r.Consumers[c.Id] = c
	return nil
}

/**
获取channel的配置
*/
func (r *RMQ) getChannelConf(configName string) (string, int, string, string, []*channel.Channel) {
	if err := apollo.InitDefault(configName); err != nil {
		log.Panicf(err.Error())
	}
	//获取基本配置
	host := agollo.GetValue("dev.host")
	port := agollo.GetIntValue("dev.port", 0)
	username := agollo.GetValue("dev.username")
	password := agollo.GetValue("dev.password")

	//获取channel的配置
	var channelList []*channel.Channel
	channelJson := agollo.GetStringValue("dev.consumers", value.EmptyString)
	json.Unmarshal([]byte(channelJson), &channelList)

	//返回基本设置和channel的list
	return host, port, username, password, channelList
}

/**
启动所有的消费者监听
*/
func RunConsumers(servername string) {
	if len(GlobalRMQ.Consumers) > 0 {
		forever := make(chan bool)
		for _, c := range GlobalRMQ.Consumers {
			if c.Name == servername {
				go func() {
					msgs := c.ConsumerMsg()
					go func() {
						for d := range msgs {
							msg := message.Message{}
							json.Unmarshal(d.Body, &msg)
							fmt.Printf("The data is %s, The msgID is %d\n", msg.Body, msg.MsgId)
						}
					}()

				}()
			}
		}
		<-forever
	}
}
