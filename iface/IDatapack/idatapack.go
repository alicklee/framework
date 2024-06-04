package IDatapack

import "github.com/CloudcadeSF/Framework/iface/message"

/**
封包、拆包模块
直接面向TCP连接中的stream，处理粘包问题
*/
type IDataPack interface {
	//获取包的头的长度
	GetHeadLen() int

	//封包
	Pack(msg message.IMessage) ([]byte, error)
	PackWithToken(msg message.IMessage, token string) ([]byte, error)
	//拆包
	Unpack([]byte) (message.IMessage, error)

	UnpackWithOutToken(binaryData []byte) (message.IMessage, error)
}
