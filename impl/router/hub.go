package router

import (
	"reflect"

	"github.com/CloudcadeSF/Framework/iface/router"
)

type Hub struct {
	apis     map[int32]router.IHandler
	apiNames map[int32]string
}

func (h *Hub) AddRouter(msgId int32, handler router.IHandler) {
	nh, ok := handler.(router.INamedHandler)
	var handlerName string
	if ok {
		handlerName = nh.GetName()
	} else {
		handlerName = reflect.ValueOf(handler).Elem().Type().String()
	}
	h.apiNames[msgId] = handlerName
	h.apis[msgId] = handler
}
func (h *Hub) GetHandlerName(msgId int32) string {
	return h.apiNames[msgId]
}
func (h *Hub) GetHandler(msgId int32) router.IHandler {
	return h.apis[msgId]
}

func NewHub() *Hub {
	return &Hub{apis: make(map[int32]router.IHandler), apiNames: map[int32]string{}}
}
