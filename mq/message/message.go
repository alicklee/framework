package message

import (
	"reflect"
	"unsafe"
)

type Message struct {
	//消息发起方
	FromPid string
	//消息接收方
	ToPid []string
	//发起方的TCP连接ID
	TcpConnId uint64
	//message的实体
	Body []byte
	//CMD
	MsgId int32
	//消息类型 0 RPC消息 1 广播消息 2 推送到个人的消息 3 推送到多个人的消息
	MsgType int
}

func (m *Message) ToByte() []byte {
	var x reflect.SliceHeader
	x.Len = int(unsafe.Sizeof(Message{}))
	x.Cap = int(unsafe.Sizeof(Message{}))
	x.Data = uintptr(unsafe.Pointer(m))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func BytesToStruct(b []byte) *Message {
	return (*Message)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}
