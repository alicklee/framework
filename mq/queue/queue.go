package queue

/**
"",    // name
false, // durable
false, // delete when unused
true,  // exclusive
false, // no-wait
nil,   // arguments
*/
type Queue struct {
	Name       string
	Durable    bool
	AutoDel    bool
	Exclusive  bool
	NoWait     bool
	RoutingKey string
	Arg        string
}

/**
创建一个新的queue
*/
func NewQueue(name string, durable bool, autoDel bool, exclusive bool, noWait bool, routingKey string, arg string) *Queue {
	q := &Queue{
		Name:       name,
		Durable:    durable,
		AutoDel:    autoDel,
		Exclusive:  exclusive,
		NoWait:     noWait,
		RoutingKey: routingKey,
		Arg:        arg,
	}
	return q
}
