package connection

import (
	"strconv"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/streadway/amqp"
)

/*
RabbitMQ 连接结构
*/
type Connection struct {
	//RabbitMQ 连接地址
	IP string
	//RabbitMQ 端口
	Port int
	//RabbitMQ 用户名
	User string
	//RabbitMQ 密码
	Password string
	//RabbitMQ 连接句柄
	Conn *amqp.Connection
}

func NewConnection(ip string, port int, user string, password string) *Connection {
	c := &Connection{
		IP:       ip,
		Port:     port,
		User:     user,
		Password: password,
	}
	return c
}

func (c *Connection) Start() (*amqp.Connection, error) {
	url := "amqp://" + c.User + ":" + c.Password + "@" + c.IP + ":" + strconv.Itoa(c.Port) + "/"
	dial, err := amqp.Dial(url)
	if err != nil {
		log.Error("Crteat mq connection error", url)
		return nil, err
	}
	return dial, nil
}

func (c *Connection) Stop() {
	c.Conn.Close()
}

func (c *Connection) Run() {
	panic("implement me")
}
