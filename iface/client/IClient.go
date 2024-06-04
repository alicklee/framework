package client

import (
	"google.golang.org/protobuf/proto"

	"github.com/CloudcadeSF/Framework/iface/response"
)

type IConnector interface {
	Connect() error
}

type ISessionFactory interface {
	NewSession(playerId string, serverId string) ISession
}

type IClient interface {
	IConnector
	ISessionFactory
}

type ISession interface {
	SendRequest(cmd uint16, reqMsg proto.Message, respMsg proto.Message) error
	SendMessage(cmd uint16, reqMsg proto.Message) error
	ProxySourceSendMessage(cmd uint16, data []byte) error
	ProxySourceRequest(cmd uint16, data []byte) response.IResponse
	SetServer(serverName, serverId string)
	SetPlayerId(playerId string)
}

type ICHashSessionFactory interface {
	NewSession(playerId string) ISession
}

type ICHashClient interface {
	IConnector
	ICHashSessionFactory
}
