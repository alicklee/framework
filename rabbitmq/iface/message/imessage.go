package message

import (
	"github.com/CloudcadeSF/Framework/iface/connection"
)

type IMessage interface {
	GetUserId() string
	GetConn() connection.IConnection
	GetBody() []byte
	Encode() []byte
	GetMsgId() uint16
}
