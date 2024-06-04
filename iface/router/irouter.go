package router

import "github.com/CloudcadeSF/Framework/iface/request"

type IRouter interface {
	//处理业务之前的hook方法
	PerHandler(r request.IRequest)
	//处理业务的hook方法
	Handler(r request.IRequest)
	//处理业务之后的hook方法
	FinishedHandler(r request.IRequest)
}
