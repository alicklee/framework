package exchange

type Exchange struct {
	//交换机的名字
	Name string
	//交换机类型
	Type string
	//消息是否持久化
	Durable bool
	//是否自动删除消息
	AutoDeleted bool
	//是否为内置交换机 设置为true的时候消费者不能往此交换机发送消息
	Internal bool
	//声明一个exchange之后是否需要服务器返回确认消息 一般设置为false不需要返回
	NoWait bool
	//额外的参数传递
	Arg string
}
