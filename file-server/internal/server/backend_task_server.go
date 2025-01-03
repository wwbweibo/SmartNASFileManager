package server

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

// BackendTaskServer used to tracking all the backend task registered to this server.
type BackendTaskServer struct {
	// tasks is a map of task name to the task
	tasks map[string]BackendTaskHandler
}

func NewBackendTaskServer() *BackendTaskServer {
	return &BackendTaskServer{
		tasks: make(map[string]BackendTaskHandler),
	}
}

func (b *BackendTaskServer) Start(ctx context.Context) error {
	var errgroup errgroup.Group
	// start a http server here
	for _, task := range b.tasks {
		log.Default().Printf("start task %s", task.GetTaskName())
		errgroup.Go(func() error {
			return task.Start(ctx)
		})
	}
	return errgroup.Wait()
}

func (b *BackendTaskServer) RegisterTask(tasks ...BackendTaskHandler) {
	for _, task := range tasks {
		b.tasks[task.GetTaskName()] = task
	}
}
