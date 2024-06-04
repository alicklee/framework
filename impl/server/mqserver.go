package server

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/CloudcadeSF/shop-heroes-legends-common/cctask"
	"github.com/CloudcadeSF/shop-heroes-legends-common/ccutil/value"
	"github.com/CloudcadeSF/shop-heroes-legends-common/id"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/CloudcadeSF/shop-heroes-legends-common/metric"
	"github.com/CloudcadeSF/shop-heroes-legends-common/network/mq"

	"github.com/CloudcadeSF/Framework/consts"
	datapack "github.com/CloudcadeSF/Framework/iface/IDatapack"
	irequest "github.com/CloudcadeSF/Framework/iface/request"
	iresponse "github.com/CloudcadeSF/Framework/iface/response"
	router2 "github.com/CloudcadeSF/Framework/iface/router"
	"github.com/CloudcadeSF/Framework/impl/connection"
	"github.com/CloudcadeSF/Framework/impl/message"
	"github.com/CloudcadeSF/Framework/impl/request"
	"github.com/CloudcadeSF/Framework/impl/router"
)

type MQServer struct {
	client     *mq.Client
	channel    *mq.Channel
	config     MQServerConfig
	routerHub  router2.IHub
	dispatcher *cctask.Dispatcher
	//msgHandler  msghandle2.IMsgHandler2
	idRouteKey  string
	allRouteKey string
	cond        *sync.Cond //初始化的同步信号
}

type MQServerConfig struct {
	//连接rabbitmq的配置
	MQUrl   string
	MQVHost string
	Name    string //服务名字
	Group   string
	Version string
	Id      string
	//RouteKey        string
	//RouteAllKey     string
	ResponseTimeout time.Duration //写入同步响应的超时时间
	DataPacker      datapack.IDataPack
}

func NewMQServer(config MQServerConfig, dispatcher *cctask.Dispatcher) *MQServer {
	s := &MQServer{
		allRouteKey: RouteKeyAll,
		config:      config,
		dispatcher:  dispatcher,
		routerHub:   router.NewHub(),
		cond:        sync.NewCond(&sync.Mutex{}),
	}
	if config.Id != "" {
		s.idRouteKey = "server:" + config.Id
	}
	return s
}

func (s *MQServer) NtfAll(data []byte) (err error) {
	return s.NtfNode("all", data)
}

func (s *MQServer) NtfNode(nodeId string, data []byte) (err error) {
	header := make(map[string]interface{})
	err = s.channel.SendNftNoCid(mq.NewSendTarget(s.exchangeName(), "server:"+nodeId, id.NextSnowFlakeInt()), data, header)
	return
}

func (s *MQServer) RequestNode(nodeId string, data []byte) (ret []byte, err error) {
	header := make(map[string]interface{})
	return s.channel.Request(mq.NewSendTarget(s.exchangeName(), "server:"+nodeId, id.NextSnowFlakeInt()), s.config.ResponseTimeout, data, header)
}

func (s *MQServer) AddHandler(cmd int32, handler router2.IHandler) {
	s.routerHub.AddRouter(cmd, handler)
}

func wrapperTaskHandler(handlerName string, handler router2.IHandler) cctask.Handler {
	return func(arg interface{}) (data interface{}, err error) {
		req := arg.(irequest.IRequest)
		log.Infof("begin do handler %s", handlerName)
		ts1 := time.Now().UnixNano()
		//handlerName := reflect.ValueOf(handler).Elem().Type().String()
		resp := handler.Handler(req)
		ts2 := time.Now().UnixNano()
		spent := float64(ts2-ts1) / float64(time.Millisecond)
		log.Infof("handler %s processed %.3f ms", handlerName, spent)
		metric.Record(handlerName, spent)
		//log.Infof("%s %s,is nil %v", handlerName, resp, resp == nil)
		if value.IsNil(resp) {
			return nil, nil
		}
		return resp, nil
	}
}

//消费收到mq包的
func (s *MQServer) consumeFunc(mgReq mq.IRequest, header map[string]interface{}) mq.IResponse {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("consumeFunc error %v stack\n%v", err, string(debug.Stack()))
		}
	}()

	//解包
	msg, err := s.config.DataPacker.UnpackWithOutToken(mgReq.GetData())
	if err != nil {
		log.Infof("mq server unpack data err %s", err.Error())
		return mq.NewResponse(nil, err)
	}
	//s.cond.L.Lock()
	//当channel 为空表示还没有初始化完成，需要等待
	//if s.channel == nil {
	//	s.cond.Wait()
	//}
	//s.cond.L.Unlock()
	//构建request
	conn := connection.NewMQConn(mgReq.GetFrom(), s.config.DataPacker, s.channel)
	if header != nil {
		pid, ok := header[consts.HeaderServerPlayerId]
		if ok {
			conn.SetPid(pid.(string))
		}
	}
	req := request.NewRequest(conn, msg, mgReq.GetCId())
	handler := s.routerHub.GetHandler(req.GetMsgId())
	if handler == nil {
		log.Errorf("handler not found msg id %d", req.GetMsgId())
		return nil
	}
	handlerName := s.routerHub.GetHandlerName(req.GetMsgId())
	var handlerRet interface{}
	var handlerErr error
	defer func() {
		if handlerErr != nil {
			//handlerName := reflect.ValueOf(handler).Elem().Type().String()
			log.Errorf("%s err %s", handlerName, handlerErr.Error())
		}
	}()

	if s.dispatcher == nil {
		handlerRet, handlerErr = wrapperTaskHandler(handlerName, handler)(req)
		if value.IsNil(handlerRet) || !mgReq.IsNeedReturn() {
			return nil
		}
	} else {
		routeKey := conn.GetPid()
		if routeKey == "" {
			routeKey = fmt.Sprintf("%d", conn.GetConnId())
		}
		//如果是异步请求就不需要返回
		if !mgReq.IsNeedReturn() {
			s.dispatcher.DispatchTask(routeKey, cctask.NewNoRetTask(req, wrapperTaskHandler(handlerName, handler)))
			return nil
		}
		//等待handler处理，并返回结果
		handlerRet, handlerErr = s.dispatcher.DoHandler(routeKey, req, wrapperTaskHandler(handlerName, handler))
		if value.IsNil(handlerRet) {
			return nil
		}
	}
	//未知的消息ID
	resp := handlerRet.(iresponse.IResponse)
	if resp.GetMsgId() == 0 {
		return nil
	}
	//打包返回数据
	respMsg := message.NewMessage(resp.GetMsgId(), resp.GetData())
	//log.Infof("mq server consumeFunc data len %d,respMsg %+v", len(respMsg.Data), respMsg)
	data, _ := s.config.DataPacker.Pack(respMsg)
	return mq.NewResponse(data, err)
}

func (s *MQServer) exchangeName() string {
	return strings.Join([]string{s.config.Name, s.config.Group, s.config.Version}, ".")
}

func (s *MQServer) Start() error {
	if s.dispatcher != nil {
		s.dispatcher.Start()
	}
	s.client = mq.NewClient(mq.ClientConfig{
		Url:   s.config.MQUrl,
		VHost: s.config.MQVHost,
	})
	routeKeys := []string{s.allRouteKey}
	if s.idRouteKey != "" {
		routeKeys = append(routeKeys, s.idRouteKey)
	}
	channelConfig := mq.ChannelConfig{
		Name:            s.config.Name,
		Exchange:        s.exchangeName(),
		ExchangeType:    "topic",
		ResponseTimeout: s.config.ResponseTimeout,
		RouteKeys:       routeKeys,
		ConsumeFunc:     s.consumeFunc,
	}
	channel, err := s.client.NewChannel(channelConfig)
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

const RouteKeyAll = "server:all"

//通知所有连接了server的客户端
func (s *MQServer) NtfMQClient(data []byte, header map[string]interface{}) {
	s.client.SendNtf(mq.NewSendTarget(s.exchangeName(), "client:all", 0), data, header)
}
