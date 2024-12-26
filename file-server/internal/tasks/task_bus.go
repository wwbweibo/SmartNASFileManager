package tasks

import (
	"fileserver/internal/server"
	"log"
)

var bus *TaskBus

func init() {
	bus = NewTaskBus()
	log.Default().Println("task bus init")
}

// TaskBus 任务总线, 用于任务之间的通信
type TaskBus struct {
	buses    map[string]chan server.ITask
	handlers map[string]server.BackendTaskHandler
}

func NewTaskBus() *TaskBus {
	return &TaskBus{
		buses:    make(map[string]chan server.ITask),
		handlers: make(map[string]server.BackendTaskHandler),
	}
}

func (b *TaskBus) Send(task server.ITask) {
	b.buses[task.GetTaskName()] <- task
}

func (b *TaskBus) RegisterHandler(task server.BackendTaskHandler) {
	b.handlers[task.GetTaskName()] = task
	b.buses[task.GetTaskName()] = make(chan server.ITask, 100)
	go func() {
		for _task := range b.buses[task.GetTaskName()] {
			task.Append(_task)
		}
	}()
}
