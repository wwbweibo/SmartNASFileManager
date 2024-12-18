package tasks

import (
	"fileserver/internal/server"
	"log"
)

var bus *TaskBus

func init() {
	bus = NewTaskBus()
	go bus.TaskHandleLoop()
	log.Default().Println("task bus init")
}

// TaskBus 任务总线, 用于任务之间的通信
type TaskBus struct {
	bus      chan server.ITask
	handlers map[string]server.BackendTaskHandler
}

func NewTaskBus() *TaskBus {
	return &TaskBus{
		bus:      make(chan server.ITask),
		handlers: make(map[string]server.BackendTaskHandler),
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
		b.handlers[task.GetTaskName()].Append(task)
	}
}
