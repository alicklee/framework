package msghandle

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/CloudcadeSF/shop-heroes-legends-common/metric"
	"github.com/CloudcadeSF/shop-heroes-legends-common/reporter"

	"github.com/CloudcadeSF/Framework/iface/request"
	"github.com/CloudcadeSF/Framework/iface/router"
	"github.com/CloudcadeSF/Framework/utils"
)

/*
message handler 的抽象类
*/
type MsgHandle struct {
	//存放每个MsgId 所对应的处理方法的map属性
	Apis map[int32]router.IRouter

	MQApis map[int32]router.IRouter

	//TCP业务工作Worker池的数量
	WorkerPoolSize uint32

	//MQ分发业务池子的数量
	MQWorkerPoolSize uint32

	//Worker负责取任务的消息队列
	TaskQueue []chan request.IRequest

	//MsgHandler的类型 0 TCP直接业务处理 1 MQ业务处理
	MsgHandlerType int
	DefaultRouter  router.IRouter
}

/**
创建一个msg handler
*/
func NewMsgHandle(handlerType int, defaultRouter router.IRouter) *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[int32]router.IRouter),
		MQApis:         make(map[int32]router.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan request.IRequest, utils.GlobalObject.WorkerPoolSize),
		//处理业务类型的type
		MsgHandlerType: handlerType,
		DefaultRouter:  defaultRouter,
	}
}

/**
直接执行TCP直连任务
*/
func (mh *MsgHandle) DoMsgHandler(request request.IRequest) {
	//防止业务函数panic导致服务器挂掉
	defer func() {
		if err := recover(); err != nil {
			var connectionId string
			msgId := request.GetMsgId()
			if request.GetConnection() != nil {
				connectionId = request.GetConnection().GetPid()
			}
			stack := debug.Stack()
			log.Errorf("player: `%s`, cmd: `%d`, get err: `%s`, stack info: `%s`", connectionId, msgId, err, stack)
			reporter.Report(fmt.Sprintf("player: `%s`, cmd: `%d`, err: `%s`,  \n, stack: %s", connectionId, request.GetMsgId(), err, stack))
		}
	}()
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		if mh.DefaultRouter != nil {
			handler = mh.DefaultRouter
		} else {
			log.Errorln("api msgId = ", request.GetMsgId(), " is not FOUND!")
			return
		}
	}
	handlerName := reflect.ValueOf(handler).Elem().Type().String()
	log.Infof("begin do handler %s", handlerName)
	//handlerName := reflect.ValueOf(handler).Elem().Type().Name()
	ts1 := time.Now().UnixNano()
	//前置处理方法
	handler.PerHandler(request)
	//处理方法
	handler.Handler(request)
	//后置处理方法
	handler.FinishedHandler(request)
	ts2 := time.Now().UnixNano()
	//log.Infof("handler %s processed %d ms", handlerName, (ts2-ts1)/int64(time.Millisecond))
	//log.Infof("handler %s processed %.3f ms", handlerName, float64(ts2-ts1)/float64(time.Millisecond))
	spent := float64(ts2-ts1) / float64(time.Millisecond)
	log.Infof("handler %s processed %.3f ms", handlerName, spent)
	metric.Record(handlerName, spent)
}

/**
收到消息之后交给对应的Api的消息类型去处理
*/
func (mh *MsgHandle) DoMQSendHandler(request request.IRequest) {
	handler, ok := mh.MQApis[request.GetMQType()]
	if !ok {
		log.Errorln("MQ api msgId = ", request.GetMQType(), " is not FOUND!")
		return
	}
	handlerName := reflect.ValueOf(handler).Elem().Type().Name()
	log.Infof("begin do handler %s", handlerName)
	ts1 := time.Now().UnixNano()
	//前置处理方法
	handler.PerHandler(request)
	//处理方法
	handler.Handler(request)
	//后置处理方法
	handler.FinishedHandler(request)
	ts2 := time.Now().UnixNano()
	//log.Infof("handler %s processed %d ms", handlerName, (ts2-ts1)/int64(time.Millisecond))
	log.Infof("handler %s processed %.3f ms", handlerName, float64(ts2-ts1)/float64(time.Millisecond))

}

/**
添加一个任务路由
*/
func (mh *MsgHandle) AddRouter(msgId int32, router router.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	//log.Infoln("Add api msgId = ", msgId)
}

func (mh *MsgHandle) AddMQRouter(mqType int32, iRouter router.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.MQApis[mqType]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(mqType)))
	}
	//2 添加msg与api的绑定关系
	mh.MQApis[mqType] = iRouter
	log.Infoln("Add MQ api MQType = ", mqType)
}

/**
发送消息到任务的消息队列中
*/
func (mh *MsgHandle) SendMsgToTaskQueue(request request.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnId() % uint64(mh.WorkerPoolSize)
	log.Info("Add ConnID=", request.GetConnection().GetConnId(), " request msgID=", request.GetMsgId(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}

/**
开启任务worker的池
*/
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//创建一个worker对应的channel的队列
		mh.TaskQueue[i] = make(chan request.IRequest, utils.GlobalObject.TaskQueueMaxLen)
		//创建一个worker
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
	log.Infof("start worker pool done ,size:%d  ", mh.WorkerPoolSize)

}

/**
创建一个worker的任务
*/
func (mh *MsgHandle) startOneWorker(id int, taskQueue chan request.IRequest) {
	//log.Infof("WorkerID : %d is started...", id)
	//开始堵塞监听对应的消息队列的消息，并交给msg handler去处理
	for {
		select {
		case data := <-taskQueue:
			mh.DoMsgHandler(data)
		}
	}
}
