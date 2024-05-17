package failed

import (
	"github.com/gopi-frame/contract/queue"
)

// MemoryFailedJob memory failed job
type MemoryFailedJob struct {
	Payload   queue.Queueable
	Queue     string
	Exception error
}

// NewMemoryFailedJob new memory failed job
func NewMemoryFailedJob(queue string, job queue.Queueable, exception error) *MemoryFailedJob {
	m := MemoryFailedJob{
		Payload:   job,
		Queue:     queue,
		Exception: exception,
	}
	return &m
}
