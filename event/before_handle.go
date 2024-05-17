package event

import "github.com/gopi-frame/contract/queue"

// NewBeforeHandle new before handle event
func NewBeforeHandle(job queue.Job) *BeforeHandle {
	return &BeforeHandle{job}
}

// BeforeHandle before handle event
type BeforeHandle struct {
	Job queue.Job
}

// Topic topic
func (event BeforeHandle) Topic() string {
	return BeforeHandleTopic
}
