package event

import "time"

// QueueAfterStop is a queue after stop event
type QueueAfterStop struct {
	Name   string
	Uptime time.Duration
}

// NewQueueAfterStop creates a new queue after stop event
func NewQueueAfterStop(name string, uptime time.Duration) *QueueAfterStop {
	return &QueueAfterStop{Name: name, Uptime: uptime}
}

// Topic returns the topic
func (QueueAfterStop) Topic() string {
	return TopicQueueAfterStop
}
