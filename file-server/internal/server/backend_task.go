package server

import (
	"context"
	"time"
)

type BackendTaskHandler interface {
	// GetTaskName return the task name
	GetTaskName() string
	// GetRunningDuration return the running duration of the task
	GetRunningDuration() time.Duration
	// Start the task,  Start method should run as a block method, it should not return until the task is finished
	Start(ctx context.Context) error
	// Stop the task
	Stop(ctx context.Context) error
	// Append a task to the task queue
	Append(task ITask)
}

type ITask interface {
	// GetTaskName return the task name
	GetTaskName() string
}
