package enum

//消息类型 0 及时回复 1 广播消息 2 推送到个人的消息 3 推送到多个人的消息
const (
	Callback = iota
	Fanout
	PersonalPush
	GroupPush
)
