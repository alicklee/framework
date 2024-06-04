package message

/**
请求消息体message的封装抽象接口
*/
type IMessage interface {
	//获取message的id
	GetMsgId() int32

	//获取message的长度
	GetMsgLen() int32

	//获取message的内容
	GetData() []byte

	//设置message的Id
	SetMsgId(id int32)

	//获取用户的Token
	GetToken() string

	//设置用户的Token
	SetToken(token string)

	//设置message的len
	SetMsgLen(n int32)

	//设置message的消息内容
	SetMsgData(data []byte)
}
