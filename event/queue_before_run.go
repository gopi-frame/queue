package event

// QueueBeforeRun is a queue before run event
type QueueBeforeRun struct {
	Name string
}

// NewQueueBeforeRun creates a new queue before run event
func NewQueueBeforeRun(name string) *QueueBeforeRun {
	return &QueueBeforeRun{Name: name}
}

// Topic returns the topic
func (QueueBeforeRun) Topic() string {
	return TopicQueueBeforeRun
}
