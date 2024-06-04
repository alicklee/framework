package impl

import (
	"encoding/json"
	"strconv"

	"github.com/CloudcadeSF/Framework/mq/channel"
	"github.com/CloudcadeSF/shop-heroes-legends-common/apollo"
	"github.com/CloudcadeSF/shop-heroes-legends-common/ccutil/value"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/streadway/amqp"
	"github.com/zouyx/agollo/v3"
)

type RabbitMQ struct {
	//RabbitMQ的链接地址
	url string
	//RabbitMQ实例的连接句柄
	Conn *amqp.Connection
	//RabbitMQ实例对应的channel
	Ch *amqp.Channel
	//RabbitMQ实例对应交换的名字
	ExchangeName string
	//RabbitMQ实例对应队列的名字
	QueueName string
	//RabbitMQ实例绑定的RoutingKey的名字
	RoutingKey string
}

/**
新建一个RabbitMQ的实例
*/
func NewRabbitMQ(configName string) *RabbitMQ {
	r := RabbitMQ{}
	host, port, username, password, _ := r.getChannelConf(configName)
	url := "amqp://" + username + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/"
	r.url = url
	dial, err := amqp.Dial(url)
	if err != nil {
		log.Panicln(err)
	}
	r.Conn = dial
	r.Ch, _ = dial.Channel()
	return &r
}

/**
启动RabbitMQ
*/
func (r *RabbitMQ) Start() error {
	return nil
}

/**
停止RabbitMQ
*/
func (r *RabbitMQ) Stop() error {
	//关闭RabbitMQ的链接
	if err := r.Conn.Close(); err != nil {
		log.Errorln("close RabbitMQ connection error", err)
		return err
	}
	//关闭RabbitMQ的channel
	if c, err := r.Conn.Channel(); err != nil {
		log.Errorln("Get RabbitMQ channel error", err)
		return err
	} else {
		if err := c.Close(); err != nil {
			log.Errorln("Close RabbitMQ channel error", err)
			return err
		}
	}
	return nil
}

/**
声明一个RabbitMQ的交换机
*/
func (r *RabbitMQ) ExchangeDeclare(exchangeName string, exchangeType string) error {
	//声明交换机
	if err := r.Ch.ExchangeDeclare(exchangeName, exchangeType, false, false, false, false, nil); err != nil {
		return err
	}
	return nil
}

/**
声明一个RabbitMQ的延时队列交换机
*/
func (r *RabbitMQ) DelayExchangeDeclare(exchangeName string, exchangeType string) error {
	args := make(amqp.Table)
	args["x-delayed-type"] = exchangeType
	if err := r.Ch.ExchangeDeclare(exchangeName, "x-delayed-message", true, false, false, false, args); err != nil {
		return err
	}
	return nil
}

/**
声明一个RabbitMQ的队列
*/
func (r *RabbitMQ) QueueDeclare(queueName string) error {
	if _, err := r.Ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}
	return nil
}

/**
绑定队列
*/
func (r *RabbitMQ) BindQueue(queueName string, exchangeName string, routingKey string) error {
	if err := r.Ch.QueueBind(queueName, routingKey, exchangeName, false, nil); err != nil {
		return err
	}
	return nil
}

/**
获取channel的配置
*/
func (r *RabbitMQ) getChannelConf(configName string) (string, int, string, string, []*channel.Channel) {
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
	if err := json.Unmarshal([]byte(channelJson), &channelList); err != nil {
		log.Errorln("Json Unmarshal data error", err)
		return "", 0, "", "", nil
	}

	//返回基本设置和channel的list
	return host, port, username, password, channelList
}
