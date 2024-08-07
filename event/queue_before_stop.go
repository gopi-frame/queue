package event

// QueueBeforeStop is a queue before stop event
type QueueBeforeStop struct {
	Name string
}

// NewQueueBeforeStop creates a new queue before stop event
func NewQueueBeforeStop(name string) *QueueBeforeStop {
	return &QueueBeforeStop{Name: name}
}

// Topic returns the topic
func (QueueBeforeStop) Topic() string {
	return TopicQueueBeforeStop
}
