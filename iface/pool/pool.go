package pool

import "github.com/CloudcadeSF/Framework/iface/task"

type IPool interface {
	Run()
	AddJob(task task.ITask)
}
