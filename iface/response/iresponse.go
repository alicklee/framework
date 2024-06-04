package response

type IResponse interface {
	//得到消息的ID
	GetMsgId() int32
	GetData() []byte
	GetError() error
}
