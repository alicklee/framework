package task

import "github.com/CloudcadeSF/Framework/iface/task"

type Task struct {
	job func() error
}

func (t *Task) Execute() {
	t.job()
}

func NewTask(f func() error) task.ITask {
	t := Task{job: f}
	return &t
}
