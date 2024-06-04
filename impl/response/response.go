package response

type Response struct {
	msgId int32
	data  []byte
	err   error
}

func (r *Response) GetData() []byte {
	return r.data
}

func (r *Response) GetError() error {
	return r.err
}

func (r *Response) GetMsgId() int32 {
	return r.msgId
}

func NewResponse(msgId int32, data []byte, err error) *Response {
	return &Response{msgId: msgId, data: data, err: err}
}
