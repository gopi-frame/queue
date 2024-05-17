package event

import (
	"github.com/gopi-frame/contract/queue"
)

// NewAfterHandle new after handle event
func NewAfterHandle(job queue.Job) *AfterHandle {
	return &AfterHandle{job}
}

// AfterHandle after handle event
type AfterHandle struct {
	Job queue.Job
}

// Topic topic
func (event AfterHandle) Topic() string {
	return AfterHandleTopic
}
