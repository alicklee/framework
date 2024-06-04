package server

import (
	"fmt"
	"net"
	"time"

	"github.com/CloudcadeSF/Framework/mq"

	connection2 "github.com/CloudcadeSF/Framework/iface/connection"

	"github.com/CloudcadeSF/Framework/impl/connection"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"

	"github.com/CloudcadeSF/Framework/iface/iconnmanager"
	"github.com/CloudcadeSF/Framework/iface/msghandle"
	"github.com/CloudcadeSF/Framework/iface/router"
	"github.com/CloudcadeSF/Framework/iface/server"
	"github.com/CloudcadeSF/Framework/impl/connmanager"
	msghandle2 "github.com/CloudcadeSF/Framework/impl/msghandle"
	"github.com/CloudcadeSF/Framework/utils"
)

type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器IP地址
	IP string
	//服务器端口号
	Port int
	//服务器版本号
	Version string
	//当前的server添加一个router
	MsgHandler msghandle.IMsgHandle
	//连接管理器
	ConnMgr iconnmanager.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn connection2.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn connection2.IConnection)
	//该Server的连接主动断开时的Hook函数
	OnConnKick func(conn connection2.IConnection)
	//通讯方式 0->TCP 1->MQ 、2->TCP和MQ都支持（用于gateway）
	MsgType int
}

var logo = `
 _______  _        _______           ______   _______  _______  ______   _______ 
(  ____ \( \      (  ___  )|\     /|(  __  \ (  ____ \(  ___  )(  __  \ (  ____ \
| (    \/| (      | (   ) || )   ( || (  \  )| (    \/| (   ) || (  \  )| (    \/
| |      | |      | |   | || |   | || |   ) || |      | (___) || |   ) || (__    
| |      | |      | |   | || |   | || |   | || |      |  ___  || |   | ||  __)   
| |      | |      | |   | || |   | || |   ) || |      | (   ) || |   ) || (      
| (____/\| (____/\| (___) || (___) || (__/  )| (____/\| )   ( || (__/  )| (____/\
(_______/(_______/(_______)(_______)(______/ (_______/|/     \|(______/ (_______/
                                                                                 
`

var top_line = `┌───────────────────────────────────────────────────┐`
var border_line = `│`
var bottom_line = `└───────────────────────────────────────────────────┘`
var copyRight = "Made By CloudCade ChengDu China Studio"

var version = "Version-0.0.1"

var github = "https://github.com/CloudcadeSF/Framework"

/**
启动服务器
*/
func (s *Server) Start() {
	go func() {
		//开启worker的pool
		s.MsgHandler.StartWorkerPool()

		//获取一个连接句柄
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Infoln("Get addr error ", err)
			return
		}
		//监听这个TCP句柄
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Infoln("listen error", err)
			return
		}
		var cid uint64
		cid = 0
		//堵塞接收客户端连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Infoln("AcceptTCP err", err)
			}
			if s.ConnMgr.Len() > utils.GlobalObject.MaxConn {
				conn.Close()
				log.Error("Max connection ")
				continue
			}
			connHandler := connection.NewConnection(s, conn, time.Duration(utils.GlobalObject.KeepAliveTime)*time.Second, cid, s.MsgHandler)
			cid++
			go connHandler.Start()
		}
	}()
}

/**
停止服务器
*/
func (s *Server) Stop() {
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

/**
运行服务器
*/
func (s *Server) Run() {
	switch s.MsgType {
	case msghandle2.MQ:
		s.StartMQModel()
	case msghandle2.TCP:
		s.Start()
		select {}
	case msghandle2.BOTH:
		s.Start()
		s.StartMQModel()
		select {}
	default:
		s.Start()
	}
}

/**
添加路由方法
*/
func (s *Server) AddRouter(msgId int32, router router.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) AddRouter2(msgId int32, handler router.IHandler) {
}

/**
添加不同MQ的路由方法
*/
func (s *Server) AddMQRouter(MQType int32, router router.IRouter) {
	s.MsgHandler.AddMQRouter(MQType, router)
}

//得到链接管理
func (s *Server) GetConnMgr() iconnmanager.IConnManager {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(connection2.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(connection2.IConnection)) {
	s.OnConnStop = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnKick(hookFunc func(connection2.IConnection)) {
	s.OnConnKick = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn connection2.IConnection) {
	if s.OnConnStart != nil {
		log.Infoln("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn connection2.IConnection) {
	if s.OnConnStop != nil {
		log.Infoln("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
	s.ConnMgr.Remove(conn)
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnKick(conn connection2.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnKick(conn)
	}
	s.ConnMgr.Remove(conn)
}

/**
开启MQ模式，开始监听MQ的消息
*/
func (s *Server) StartMQModel() {
	mq.Start("rabbitmq", s.Name)
}

func (s *Server) OnInit() {

}

/**
销毁服务
*/
func (s *Server) OnDestroy() {
	log.Infof("network server: `%s` will exit...", s.Name)
	// 清理连接管理器
	s.ConnMgr.ClearConn()
	log.Infof("network server: `%s` exit...", s.Name)
}

/**
实例化一个服务器
*/
func InitTcpServer(name string, ipVersion string, ip string, port int, version string, defaultRouter router.IRouter) server.IServer {
	log.Infof("%s(version: %s) is starting ip version is %s ip is %s port is %d",
		name,
		version,
		ipVersion,
		ip,
		port)
	s := &Server{
		Name:       name,
		IPVersion:  ipVersion,
		IP:         ip,
		Port:       port,
		Version:    version,
		MsgHandler: msghandle2.NewMsgHandle(msghandle2.TCP, defaultRouter),
		ConnMgr:    connmanager.NewConnManager(),
		MsgType:    msghandle2.TCP,
	}
	return s
}

func init() {
	fmt.Println(logo)
	fmt.Println(top_line)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/CloudcadeSF/Framework %s", border_line, border_line))
	fmt.Println(fmt.Sprintf("%s    Made by CloudCade ChengDu China studio         %s", border_line, border_line))
	fmt.Println(fmt.Sprintf("%s                 %s                     %s", border_line, version, border_line))
	fmt.Println(bottom_line)
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
