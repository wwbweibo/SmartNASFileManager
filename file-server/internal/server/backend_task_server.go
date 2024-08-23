package server

import "context"

// BackendTaskServer used to tracking all the backend task registered to this server.
type BackendTaskServer struct {
	// tasks is a map of task name to the task
	tasks map[string]BackendTask
}

func NewBackendTaskServer() *BackendTaskServer {
	return &BackendTaskServer{}
}

func (b *BackendTaskServer) Start(ctx context.Context) error {
	// start a http server here
	return nil
}

func (b *BackendTaskServer) RegisterTask(tasks ...BackendTask) {
	for _, task := range tasks {
		b.tasks[task.GetTaskName()] = task
	}
}
