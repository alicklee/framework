package property

/**
连接属性结构体
*/
type Property struct {
	Pid   string
	Token string
}

/**
设置一个property
*/
func (p *Property) SetProperty(pid string, token string) {
	p.Pid = pid
	p.Token = token
}
