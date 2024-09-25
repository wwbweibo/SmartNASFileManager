package server

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// BackendTaskServer used to tracking all the backend task registered to this server.
type BackendTaskServer struct {
	// tasks is a map of task name to the task
	tasks map[string]BackendTask
}

func NewBackendTaskServer() *BackendTaskServer {
	return &BackendTaskServer{
		tasks: make(map[string]BackendTask),
	}
}

func (b *BackendTaskServer) Start(ctx context.Context) error {
	var errgroup errgroup.Group
	// start a http server here
	for _, task := range b.tasks {
		errgroup.Go(func() error {
			return task.Start(ctx)
		})
	}
	return errgroup.Wait()
}

func (b *BackendTaskServer) RegisterTask(tasks ...BackendTask) {
	for _, task := range tasks {
		b.tasks[task.GetTaskName()] = task
	}
}
