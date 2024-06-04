package message

import (
	"bytes"
	"encoding/gob"
	"github.com/CloudcadeSF/Framework/iface/connection"
	"log"
)

type Message struct {
	Conn connection.IConnection
	Body []byte
	FromServerId string
	ToServerId   string
	UserId string
	MsgId uint16
}

func init() {
	gob.Register(Message{})
}

func NewMessage(conn connection.IConnection,body []byte,userId string) *Message{
	m := & Message{
		Conn: conn,
		Body: body,
		UserId: userId,
	}
	return m
}

func (m *Message) GetUserId() string {
	return m.UserId
}

func (m *Message) GetMsgId() uint16 {
	return m.MsgId
}

func (m *Message) GetConn() connection.IConnection {
	return m.Conn
}

func (m *Message) GetBody() []byte {
	return m.Body
}

func (m *Message) Encode() []byte {
	 b, err := encode(m)
	if err != nil {
		return nil
	}
	return b
}

func Decode(b []byte) *Message {
	message, err := decode(b)
	if err != nil{
		return nil
	}
	return &message
}

// 编码，把结构体数据编码成字节流数据
func encode(m *Message) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf) // 构造编码器，并把数据写进buf中
	if err := encoder.Encode(m); err != nil {
		log.Printf("encode error: %v\n", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// 解码，把字节流数据解析成结构体数据
func decode(b []byte) (Message, error) {
	bufPtr := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(bufPtr)
	var p Message
	if err := decoder.Decode(&p); err != nil {
		log.Println("Decode RabbitMQ message error",err)
		return Message{}, err
	}
	return p, nil
}


