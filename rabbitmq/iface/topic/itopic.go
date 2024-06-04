package topic

/**
Topic交换机接口
*/
type ITopic interface {
	//发送消息
	SendMsg(msg []byte, routingKey string) error
	SendListMsg(msg []byte,pids []string) error
	FanoutMsg(msg []byte) error
	//接受消息
	RecvMsg() error
	//接受消息并转发
	RecvAndRelay() error
	//开启Topic
	Start() error
	//关闭Topic
	Stop()
}
