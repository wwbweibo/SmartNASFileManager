package tasks

import (
	"fileserver/internal/server"
)

var bus *TaskBus

func init() {
	bus = NewTaskBus()
}

// TaskBus 任务总线, 用于任务之间的通信
type TaskBus struct {
	bus      chan server.ITask
	handlers map[string]server.BackendTaskHandler
}

func NewTaskBus() *TaskBus {
	return &TaskBus{
		bus: make(chan server.ITask),
	}
}

func (b *TaskBus) Send(task server.ITask) {
	b.bus <- task
}

func (b *TaskBus) RegisterHandler(task server.BackendTaskHandler) {
	b.handlers[task.GetTaskName()] = task
}

func (b *TaskBus) TaskHandleLoop() {
	for task := range b.bus {
		// do something
		b.handlers[task.GetTaskName()].Append(task)
	}
}
