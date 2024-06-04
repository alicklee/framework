package connection

import (
	"net"
	"time"

	"github.com/CloudcadeSF/Framework/iface/message"
	"github.com/CloudcadeSF/Framework/iface/response"
)

type IConnection interface {
	Start()
	Stop()
	Kick()
	GetTCPConnection() *net.TCPConn
	GetConnId() uint64
	RemoteAddr() net.Addr
	Request(msgId int32, data []byte, timeout time.Duration) response.IResponse
	SendMsg(msgId int32, data []byte) error
	SendMsgImmediately(msgId int32, data []byte) error
	SendProxy(data []byte) error
	SendMsgList(messages ...message.IMessage) error
	FanoutMsg(msgId int32, data []byte) error
	GetPid() string
	SetPid(pid string)

	//cid 透传使用 参数中cid不为空时，会将cid传出。
	//为空时会自动填写cid，并返回其值
	SendMsgWithCID(cid string, msgId int32, data []byte) (string, error)
}

type Header map[string]interface{}

type HandleFunc func(*net.TCPConn, []byte, int) error
