package connmanager

import (
	"errors"
	"sync"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/CloudcadeSF/shop-heroes-legends-common/metric"

	"github.com/CloudcadeSF/Framework/iface/connection"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint64]connection.IConnection //管理的连接信息
	connLock    sync.RWMutex                      //读写连接的读写锁
}

/*
	创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]connection.IConnection),
	}
}

/**
添加一个连接管理
*/
func (connMgr *ConnManager) Add(conn connection.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	metric.SetOnlineUser(connMgr.Len())
	//将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnId()] = conn

	log.Infoln("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

/**
移除一个连接
*/
func (connMgr *ConnManager) Remove(conn connection.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	metric.SetOnlineUser(connMgr.Len())
	//删除连接信息
	delete(connMgr.connections, conn.GetConnId())
}

/**
利用ConnID获取链接
*/
func (connMgr *ConnManager) Get(connID uint64) (connection.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

/**
获取一个pid slice里面所有的连接
*/
func (connMgr *ConnManager) FindConnsByPids(pids []string) []connection.IConnection {
	var result []connection.IConnection
	for _, v := range connMgr.connections {
		for _, pid := range pids {
			if v.GetPid() == pid {
				result = append(result, v)
			}
		}
	}
	return result
}

/**
获取对应Pid的连接
*/
func (connMgr *ConnManager) FindConnByPid(pid string) connection.IConnection {
	for _, v := range connMgr.connections {
		if v.GetPid() == pid {
			return v
		}
	}
	return nil
}

func (connMgr *ConnManager) FindAll() []connection.IConnection {
	var result []connection.IConnection
	for _, conn := range connMgr.connections {
		result = append(result, conn)
	}
	return result
}

/**
获取当前连接数
*/
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

/**
清楚所有的连接信息
*/
func (connMgr *ConnManager) ClearConn() {
	var total int32
	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Kick()
		//删除
		delete(connMgr.connections, connID)
		total++
	}
	log.Errorf("clear all connection successfully: removed count: `%d`", total)
}
