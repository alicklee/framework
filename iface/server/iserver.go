package server

import (
	connection2 "github.com/CloudcadeSF/Framework/iface/connection"
	"github.com/CloudcadeSF/Framework/iface/iconnmanager"
	"github.com/CloudcadeSF/Framework/iface/router"
)

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	OnDestroy()
	OnInit()
	//运行服务器
	Run()
	//添加一个路由
	AddRouter(msgId int32, router router.IRouter)
	//添加一个路由
	AddRouter2(msgId int32, handler router.IHandler)

	//添加一个MQ的路由(router router.IRouter)
	AddMQRouter(MQtype int32, router router.IRouter)

	//设置该Server的连接创建时Hook函数
	SetOnConnStart(hookFunc func(connection2.IConnection))

	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(hookFunc func(connection2.IConnection))

	//设置该Server的连接主动断开时的Hook函数
	SetOnConnKick(hookFunc func(connection2.IConnection))

	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn connection2.IConnection)

	//调用连接OnConnKick Hook函数
	CallOnConnKick(conn connection2.IConnection)

	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn connection2.IConnection)

	//获取连接管理
	GetConnMgr() iconnmanager.IConnManager

	//开启MQ模式
	StartMQModel()
}
