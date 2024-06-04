package router

import "github.com/CloudcadeSF/Framework/iface/request"

/**
这里做了一个抽象的封装，实现router的时候，先嵌入这个基类，才根据业务在重写对应的方法
*/
type BasicRouter struct{}

//处理业务之前的hook方法
func (br *BasicRouter) PerHandler(r request.IRequest) {}

//处理业务的hook方法
func (br *BasicRouter) Handler(r request.IRequest) {}

//处理业务之后的hook方法
func (br *BasicRouter) FinishedHandler(r request.IRequest) {}


