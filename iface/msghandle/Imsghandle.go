package msghandle

import (
	"github.com/CloudcadeSF/Framework/iface/request"
	"github.com/CloudcadeSF/Framework/iface/router"
)

/*
消息管理抽象层
*/
type IMsgHandle interface {
	DoMsgHandler(request request.IRequest)            //马上以非阻塞方式处理消息
	AddRouter(msgId int32, iRouter router.IRouter)    //为消息添加具体的处理逻辑
	AddMQRouter(mqType int32, iRouter router.IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                                 //启动worker工作池
	SendMsgToTaskQueue(request request.IRequest)      //将消息交给TaskQueue,由worker进行处理
}
