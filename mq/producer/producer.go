package producer

import (
	"encoding/json"

	"github.com/CloudcadeSF/Framework/iface/request"
	"github.com/CloudcadeSF/Framework/mq/message"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"

	"github.com/CloudcadeSF/Framework/mq/channel"

	"github.com/streadway/amqp"
)

var GlobalProducer map[int]*Producer = make(map[int]*Producer)

/**
生产者结构
*/
type Producer struct {
	//生产者通道
	Name string
	//Channel
	Ch *channel.Channel
	//生产者编号
	Id int
}

/**
创建一个生产者实体
*/
func NewProducer(id int, name string, ch *channel.Channel) *Producer {
	p := Producer{
		Name: name,
		Ch:   ch,
		Id:   id,
	}
	return &p
}

/**
发送消息到MQ
*/
func (p *Producer) SendMsg(req request.IRequest) error {
	m := message.Message{
		TcpConnId: req.GetConnection().GetConnId(),
		Body:      req.GetData(),
		FromPid:   req.GetConnection().GetPid(),
		MsgId:     req.GetMsgId(),
	}
	marshal, _ := json.Marshal(&m)
	ch := p.Ch.Ch
	err := ch.Publish(
		p.Ch.Exchange.Name,    // exchange
		p.Ch.Queue.RoutingKey, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        marshal,
		})
	if err != nil {
		return err
	}
	log.Infof("Producter id is %d Name is %s sent the msg %s", p.Id, p.Name, m.Body)
	return nil
}
