package request

import (
	"github.com/CloudcadeSF/Framework/iface/connection"
	"github.com/CloudcadeSF/Framework/iface/message"
)

type Request struct {
	conn   connection.IConnection
	msg    message.IMessage
	cid    string
	MQType int32
}

func (r *Request) GetConnection() connection.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
func (r *Request) GetMsgId() int32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetMQType() int32 {
	return r.MQType
}

func (r *Request) SetMQType(MQType int32) {
	r.MQType = MQType
}

func (r *Request) GetToken() string {
	return r.msg.GetToken()
}

func (r *Request) GetCId() string {
	return r.cid
}

func NewRequest(c connection.IConnection, msg message.IMessage, cid string) *Request {
	req := Request{
		conn: c,
		msg:  msg,
		cid:  cid,
	}
	return &req
}
