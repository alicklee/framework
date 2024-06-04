package request

import "github.com/CloudcadeSF/Framework/iface/connection"

type IRequest interface {
	//等到request的连接
	GetConnection() connection.IConnection
	//得到消息
	GetData() []byte
	//得到消息的ID
	GetMsgId() int32
	//获取MQ的类型
	GetMQType() int32
	//设置MQ的类型
	SetMQType(MQType int32)
	//获取消息中的token
	GetToken() string
	//获取消息的CID
	GetCId() string
}
