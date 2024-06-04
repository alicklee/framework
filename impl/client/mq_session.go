package client

import (
	"time"

	"github.com/CloudcadeSF/Framework/iface/response"
	"github.com/CloudcadeSF/Framework/impl/connection"
	response2 "github.com/CloudcadeSF/Framework/impl/response"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"google.golang.org/protobuf/proto"
)

type MQSession struct {
	requestTimeout time.Duration
	conn           *connection.MQConnection
}

func (s *MQSession) SendRequest(cmd int32, reqMsg proto.Message, respMsg proto.Message) error {
	bytes, err := proto.Marshal(reqMsg)
	if err != nil {
		log.Infof("marshal data error %s", err.Error())
		return err
	}
	resp := s.conn.Request(cmd, bytes, s.requestTimeout)
	if resp == nil {
		return nil
	}
	if err := resp.GetError(); err != nil {
		return err
	}
	return proto.Unmarshal(resp.GetData(), respMsg)
}

func (s *MQSession) SendMessage(cmd int32, reqMsg proto.Message) error {
	bytes, err := proto.Marshal(reqMsg)
	if err != nil {
		log.Infof("marshal data error %s", err.Error())
		return err
	}
	return s.conn.SendMsg(cmd, bytes)
}

func (s *MQSession) ProxySourceSendMessage(cid string, cmd int32, data []byte) (string, error) {
	return s.conn.SendMsgWithCID(cid, cmd, data)
}

func (s *MQSession) ProxySourceRequest(cmd int32, data []byte) response.IResponse {
	return s.conn.Request(cmd, data, s.requestTimeout)
}

//透明转发
func (s *MQSession) ProxyRequest(cmd int32, reqMsg proto.Message) response.IResponse {
	data, err := proto.Marshal(reqMsg)
	if err != nil {
		log.Infof("marshal data error %s", err.Error())
		return response2.NewResponse(0, nil, err)
	}
	return s.ProxySourceRequest(cmd, data)
}

func (s *MQSession) SetServer(serverName, serverId string) {
	s.conn.Target.Exchange = serverName
	s.conn.Target.RouteKey = "server:" + serverId
}

func (s *MQSession) SetPlayerId(playerId string) {
	s.conn.SetPid(playerId)
}
