package pool

import (
	"github.com/CloudcadeSF/Framework/iface/pool"
	"github.com/CloudcadeSF/Framework/iface/task"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
)

type Pool struct {
	EntryChannel chan task.ITask
	JobsChannel  chan task.ITask
	num          int
}

func NewPool(n int) pool.IPool {
	p := &Pool{
		EntryChannel: make(chan task.ITask),
		JobsChannel:  make(chan task.ITask),
		num:          n,
	}
	return p
}

func (p *Pool) AddJob(task task.ITask) {
	p.EntryChannel <- task
}

func (p *Pool) Worker(id int) {
	for task := range p.JobsChannel {
		task.Execute()
		log.Infoln("Worker id ： ", id, " executed")
	}
}

func (p *Pool) Run() {

	for i := 0; i < p.num; i++ {
		go p.Worker(i)
	}

	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}
}
