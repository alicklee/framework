package task

import "github.com/CloudcadeSF/Framework/iface/task"

type Task2 struct {
	job func() error
}

func (n *Task2) Execute() {
	n.job()
}

func NewTask2(f func() error) task.ITask {
	t := Task2{job: f}
	return &t
}
