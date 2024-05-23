package memory

import "github.com/gopi-frame/contract/queue"

// FailedJob failed job
type FailedJob struct {
	Payload   queue.Queueable
	Queue     string
	Exception error
}
