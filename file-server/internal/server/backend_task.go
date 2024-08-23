package server

import (
	"context"
	"time"
)

type BackendTask interface {
	// GetTaskName return the task name
	GetTaskName() string
	// GetRunningDuration return the running duration of the task
	GetRunningDuration() time.Duration
	// Start the task
	Start(ctx context.Context) error
	// Stop the task
	Stop(ctx context.Context) error
}
