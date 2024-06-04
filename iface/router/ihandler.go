package router

import (
	"github.com/CloudcadeSF/Framework/iface/request"
	"github.com/CloudcadeSF/Framework/iface/response"
)

type IHandler interface {
	//处理业务的hook方法
	Handler(request.IRequest) response.IResponse
}

type INamedHandler interface {
	GetName() string
}
