package delay

type IDelay interface {
	//发送消息
	SendMsg([]byte,int64) error
	//接受消息
	RecvMsg(msgHandler func(interface{})) error
	//接受消息并转发
	RecvAndRelay() error
	//设施延时队列的header
	SetHeader() error
	//开启Topic
	Start() error
	//关闭Topic
	Stop()
}
