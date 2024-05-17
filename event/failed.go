package event

import "github.com/gopi-frame/contract/queue"

// Failed failed event
type Failed struct {
	Job   queue.Job
	Error error
}

// NewFailed new failed handle event
func NewFailed(job queue.Job, err error) *Failed {
	return &Failed{job, err}
}

// Topic return topic
func (event Failed) Topic() string {
	return FailedTopic
}
