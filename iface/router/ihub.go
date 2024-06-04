package router

type IHub interface {
	AddRouter(msgId int32, handler IHandler)
	GetHandler(msgId int32) IHandler
	GetHandlerName(msgId int32) string
}
