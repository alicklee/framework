package message

/**
定义通讯消息体的结构
*/
type Message struct {
	//消息的ID（占4字节）
	Id int32
	//消息的长度（占4字节）
	DataLen int32
	//用户的token(占64字节)
	Token string
	//消息的内容
	Data []byte
}

func NewMessage(msgId int32, data []byte) *Message {
	msg := &Message{
		Id:      msgId,
		DataLen: int32(len(data)),
		Data:    data,
	}
	return msg
}

//获取消息的ID
func (m *Message) GetMsgId() int32 {
	return m.Id
}

//获取消息的长度
func (m *Message) GetMsgLen() int32 {
	return m.DataLen
}

//获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息的ID
func (m *Message) SetMsgId(id int32) {
	m.Id = id
}

//获取消息的token
func (m *Message) GetToken() string {
	return m.Token
}

//设置消息的token
func (m *Message) SetToken(token string) {
	m.Token = token
}

//设置消息的长度
func (m *Message) SetMsgLen(n int32) {
	m.DataLen = n
}

//设置消息的内容
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}
