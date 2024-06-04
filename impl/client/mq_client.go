package client

import (
	"strings"
	"time"

	"github.com/CloudcadeSF/shop-heroes-legends-common/ccutil/value"
	"github.com/CloudcadeSF/shop-heroes-legends-common/id"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/CloudcadeSF/shop-heroes-legends-common/network/mq"

	"github.com/CloudcadeSF/Framework/consts"
	datapack "github.com/CloudcadeSF/Framework/iface/IDatapack"
	"github.com/CloudcadeSF/Framework/iface/message"
	"github.com/CloudcadeSF/Framework/impl/connection"
	datapack2 "github.com/CloudcadeSF/Framework/impl/datapack"
)

type MQClient struct {
	client      *mq.Client
	channel     *mq.Channel
	config      MQClientConfig
	idRouteKey  string
	allRouteKey string
	//cond        *sync.Cond //初始化的同步信号
}

type MQClientConfig struct {
	//连接rabbitmq的配置
	ServerMQUrl     string
	ServerMQVHost   string
	ServerName      string
	ServerGroup     string
	ServerVersion   string
	Id              string        //客户端ID，用来区分同一server下的其他客户端
	ResponseTimeout time.Duration //写入同步响应的超时时间
	DataPacker      datapack.IDataPack
	ConsumeFunc     MessageConsumeFunc //收到消息后的处理函数
}
type MessageConsumeFunc func(playerId string, req mq.IRequest, msg message.IMessage) (message.IMessage, error)

func defaultMessageConsumeFunc(playerId string, req mq.IRequest, msg message.IMessage) (message.IMessage, error) {
	return nil, nil
}

const RouteKeyAll = "client:all"

func NewMQClient(config MQClientConfig) *MQClient {
	if config.ConsumeFunc == nil {
		config.ConsumeFunc = defaultMessageConsumeFunc
	}
	if config.DataPacker == nil {
		config.DataPacker = datapack2.SimpleDataPacker
	}
	c := &MQClient{
		config:      config,
		allRouteKey: RouteKeyAll,
		//cond:        sync.NewCond(&sync.Mutex{}),
	}
	if config.Id != "" {
		c.idRouteKey = "client:" + config.Id
	}
	//log.Infof("NewMQClient %+v", c)
	return c
}
func (s *MQClient) Close() error {
	return s.client.Close()
}

func (s *MQClient) mqConnConsumeFunc(mgReq mq.IRequest, header map[string]interface{}) mq.IResponse {
	//s.cond.L.Lock()
	//当channel 为空表示还没有初始化完成，需要等待
	//if s.channel == nil {
	//	s.cond.Wait()
	//}
	//s.cond.L.Unlock()
	var pid = ""
	if header != nil {
		pidV, ok := header[consts.HeaderServerPlayerId]
		if ok {
			pid = pidV.(string)
		}
	}
	msg, err := s.config.DataPacker.UnpackWithOutToken(mgReq.GetData())
	if err != nil {
		log.Infof("mq server unpack data err %s", err.Error())
		return mq.NewResponse(nil, err)
	}
	respMsg, err := s.config.ConsumeFunc(pid, mgReq, msg)
	if err != nil {
		return mq.NewResponse(nil, err)
	}
	if respMsg == nil {
		return nil
	}
	data, err := s.config.DataPacker.Pack(respMsg)
	return mq.NewResponse(data, err)
}

func (s *MQClient) exchangeName() string {
	return strings.Join([]string{s.config.ServerName, s.config.ServerGroup, s.config.ServerVersion}, ".")
}

func (s *MQClient) Connect() error {
	s.client = mq.NewClient(mq.ClientConfig{
		Url:   s.config.ServerMQUrl,
		VHost: s.config.ServerMQVHost,
	})
	//if err := s.client.Connect(); err != nil {
	//	return err
	//}
	routeKeys := []string{s.allRouteKey}
	if s.idRouteKey != "" {
		routeKeys = []string{s.idRouteKey, s.allRouteKey} //append(routeKeys, s.idRouteKey)
	}
	channel, err := s.client.NewChannel(mq.ChannelConfig{
		Name:            value.RandString(4),
		Exchange:        s.exchangeName(),
		ExchangeType:    "topic",
		ResponseTimeout: s.config.ResponseTimeout,
		ConsumeFunc:     s.mqConnConsumeFunc,
		RouteKeys:       routeKeys,
	})
	if err != nil {
		return err
	}
	//s.cond.L.Lock()
	s.channel = channel
	//初始化完成，通知consumer开始工作
	//s.cond.Broadcast()
	//s.cond.L.Unlock()
	return s.client.Connect()
}

func (s *MQClient) NewSession(playerId string, serverId string) *MQSession {
	if serverId == "" {
		serverId = "all"
	}
	conn := connection.NewMQConn(
		mq.NewSendTarget(s.exchangeName(), "server:"+serverId, id.NextSnowFlakeInt()),
		s.config.DataPacker, s.channel)
	conn.SetPid(playerId)
	return &MQSession{
		requestTimeout: s.config.ResponseTimeout,
		conn:           conn,
	}
}

func (s *MQClient) NewAllSession(playerId string) *MQSession {
	return s.NewSession(playerId, "")
}
