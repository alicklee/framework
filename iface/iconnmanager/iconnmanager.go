package iconnmanager

import "github.com/CloudcadeSF/Framework/iface/connection"

/*
	连接管理抽象层
*/
type IConnManager interface {
	//添加链接
	Add(conn connection.IConnection)
	//删除连接
	Remove(conn connection.IConnection)
	//利用ConnID获取链接
	Get(connID uint64) (connection.IConnection, error)
	//获取当前连接数量
	Len() int
	//删除并停止所有链接
	ClearConn()
	//获取一个pid slice里面所有的连接
	FindConnsByPids(pids []string) []connection.IConnection
	//获取对应Pid的连接
	FindConnByPid(pid string) connection.IConnection
	//查找全部链接
	FindAll() []connection.IConnection
}
