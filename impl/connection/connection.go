package connection

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	connection2 "github.com/CloudcadeSF/Framework/iface/connection"
	"github.com/CloudcadeSF/Framework/iface/iconnmanager"
	"github.com/CloudcadeSF/Framework/iface/response"

	"github.com/CloudcadeSF/Framework/impl/request"

	message2 "github.com/CloudcadeSF/Framework/impl/message"

	"github.com/CloudcadeSF/Framework/iface/message"
	"github.com/CloudcadeSF/Framework/iface/server"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"

	"github.com/CloudcadeSF/Framework/iface/msghandle"
	"github.com/CloudcadeSF/Framework/impl/datapack"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer server.IServer
	//socket TCP
	Conn *net.TCPConn
	//链接的ID
	ConnID uint64
	//链接是否被关闭
	IsClose bool
	//等待链接关闭的channel
	ExitConnChan chan bool
	//添加message的channel
	MsgChan       chan []byte
	KeepAliveTime time.Duration
	//该链接处理的方法router
	MsgHandler msghandle.IMsgHandle
	//用户ID
	Pid string

	//只执行一次stop操作
	terminateOnce sync.Once
}

func NewConnection(server server.IServer, conn *net.TCPConn, keepAliveTime time.Duration, connId uint64, msgHandler msghandle.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:     server,
		Conn:          conn,
		ConnID:        connId,
		IsClose:       false,
		KeepAliveTime: keepAliveTime,
		MsgHandler:    msgHandler,
		MsgChan:       make(chan []byte),
		ExitConnChan:  make(chan bool, 1),
	}
	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) Start() {
	log.Infof("[framework tcp] connection connId %v started", c.ConnID)
	//把连接加入到连接管理模块
	//c.ConnManager.Add(c)
	//启动从当前链接读取数据的业务
	go c.startReader()
	//启动从当前连接回写数据的业务
	go c.StartWrite()
	//执行创建连接的hook函数
	c.TcpServer.CallOnConnStart(c)

}

func (c *Connection) terminate(cb func(conn connection2.IConnection)) {
	if c.IsClose == true {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("connection.terminate panic: %v", err)
		}

		c.IsClose = true
		//告知Reader已经关闭连接
		c.ExitConnChan <- true

		//回收管道资源
		close(c.ExitConnChan)
		close(c.MsgChan)

		//关闭连接
		if err := c.Conn.Close(); err != nil {
			log.Panicln("Connection close error")
		}
	}()

	if cb == nil {
		return
	}

	cb(c)
}

// 关闭一个连接
func (c *Connection) Stop() {
	c.terminateOnce.Do(func() {
		log.Infof("[framework tcp] connection connId %v stop by service", c.ConnID)
		c.terminate(c.TcpServer.CallOnConnStop)
	})
}

// 关闭一个连接
func (c *Connection) Kick() {
	c.terminateOnce.Do(func() {
		log.Infof("[framework tcp] connection connId %v stop by kick", c.ConnID)
		c.terminate(c.TcpServer.CallOnConnKick)
	})
}

func (c *Connection) Request(msgId int32, data []byte, timeout time.Duration) response.IResponse {
	panic("not implement")
}

// 获取一个连接
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

/*
*
获取连接的ID
*/
func (c *Connection) GetConnId() uint64 {
	return c.ConnID
}

/*
*
获取客户端的addr
*/
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

/*
启动写携程
*/
func (c *Connection) StartWrite() {
	log.Infof("[framework tcp] connection connId %v writer started", c.ConnID)
	defer log.Infof("[framework tcp] connection connId %v writer stoped", c.ConnID)

	for {
		select {
		case data := <-c.MsgChan:
			//接收到channel里面的消息就开始回写消息给客户端
			if _, err := c.Conn.Write(data); err != nil {
				log.Errorf("[framework tcp] connection connId %v write failed %v",
					c.ConnID, err)
				c.Stop()
				return
			}
			//g, _ := c.TcpServer.GetConnMgr().Get(c.ConnID)
			//log.Info("manager connection pid is ", g.GetPid())
		case <-c.ExitConnChan:
			//当用户退出了，关闭连接
			log.Infof("[framework tcp] connection connId %v writer close", c.ConnID)
			return
		}
	}
}

/*
*
启动一个reader开始读取连接的数据
*/
func (c *Connection) startReader() {
	log.Infof("[framework tcp] connection connId %v reader start", c.ConnID)
	//读取完毕断开连接
	defer func() {
		c.Stop()
		log.Infof("[framework tcp] connection connId %v reader stoped", c.ConnID)
	}()

	//堵塞开始读取数据
	for {
		dp := datapack.NewDataPack()
		msg, err := c.unpack(dp)
		if err != nil {
			log.Errorf("[framework tcp] connection connId %v reader err %v",
				c.ConnID, err)
			break
		}
		//执行注册的路由方法
		req := request.NewRequest(c, msg, "unsupported")

		//把消息发送给工作池的消息队列
		c.MsgHandler.SendMsgToTaskQueue(req)
	}
}

/*
*
解包为message结构的方法
*/
func (c *Connection) unpack(dp *datapack.DataPack) (message.IMessage, error) {
	//读取客户端消息的head，lenght + i
	headData := make([]byte, dp.GetHeadLen())
	tcpConn := c.GetTCPConnection()

	//读取超时时间
	if err := tcpConn.SetKeepAlivePeriod(c.KeepAliveTime); err != nil {
		return nil, fmt.Errorf("read timeout wait %v err %v",
			c.KeepAliveTime, err)
	}

	if _, err := io.ReadFull(tcpConn, headData); err != nil {
		return nil, fmt.Errorf("read package header len %v err %v",
			len(headData), err)
	}

	//拆包，得到msgId 和 msgLen 和 token
	msg, err := dp.Unpack(headData)
	if err != nil {
		return nil, fmt.Errorf("unpack header headData %+v err %v",
			headData, err)
	}

	//根据dataLen 再次读取data
	var data []byte
	if msg.GetMsgLen() > 0 {
		data = make([]byte, msg.GetMsgLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
			return nil, fmt.Errorf("read msgId %v body len %v headerData %+v err %v",
				msg.GetMsgId(), msg.GetMsgLen(), headData, err)
		}
	}

	//填充message结构体
	msg.SetMsgData(data)
	return msg, nil
}

/*
*
发送消息到当前的connId
*/
func (c *Connection) GetPid() string {
	return c.Pid
}

func (c *Connection) SetPid(pid string) {
	c.Pid = pid
}

func (c *Connection) SendProxy(data []byte) error {
	_, err := c.GetTCPConnection().Write(data)
	return err
}

func (c *Connection) SendMsgImmediately(msgId int32, data []byte) error {
	if c.IsClose {
		return errors.New("Connection has closed")
	}
	dp := datapack.NewDataPack()
	msg := message2.NewMessage(msgId, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		log.Infoln("pack error message id = ", msg.Id)
		return errors.New("pack message error")
	}
	_, err = c.GetTCPConnection().Write(binaryMsg)
	return err
}

/*
*
发送消息
*/
func (c *Connection) SendMsg(msgId int32, data []byte) error {
	if c.IsClose {
		return errors.New("Connection has closed")
	}
	dp := datapack.NewDataPack()
	msg := message2.NewMessage(msgId, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		log.Infoln("pack error message id = ", msg.Id)
		return errors.New("pack message error")
	}
	//将消息写入无缓冲的通道，交给Write去处理
	c.MsgChan <- binaryMsg
	return nil
}

func (c *Connection) SendMsgWithCID(cid string, mid int32, data []byte) (string, error) {
	return "", c.SendMsg(mid, data)
}

/*
*
选定player的广播
*/
func (c *Connection) SendMsgList(messages ...message.IMessage) error {
	for _, iMessage := range messages {
		if err := c.SendMsg(iMessage.GetMsgId(), iMessage.GetData()); err != nil {
			return err
		}
	}
	return nil
	//var connM iconnmanager.IConnManager
	//if conns := connM.FindConnsByPids(pidList); conns != nil {
	//	for _, conn := range conns {
	//		if err := conn.SendMsg(msgId, data); err != nil {
	//			log.Errorln("Send msg error", err)
	//			return err
	//		}
	//	}
	//}
	//return nil
}

/*
*
全服广播
*/
func (c *Connection) FanoutMsg(msgId int32, data []byte) error {
	var connM iconnmanager.IConnManager
	conns := connM.FindAll()
	if len(conns) > 0 {
		for _, conn := range conns {
			if err := conn.SendMsg(msgId, data); err != nil {
				log.Errorln("Send msg error", err)
				return err
			}
		}
	}
	return nil
}
