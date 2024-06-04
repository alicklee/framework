package producer

type IProducer interface {
	//发送消息
	SendMsg(msg []byte) error
}
