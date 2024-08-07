package event

// JobAfterHandle is a job after handle event
type JobAfterHandle struct {
	Name string
	ID   string
}

// NewJobAfterHandle creates a new job after handle event
func NewJobAfterHandle(id string, name string) *JobAfterHandle {
	return &JobAfterHandle{ID: id, Name: name}
}

// Topic returns the topic
func (JobAfterHandle) Topic() string {
	return TopicJobAfterHandle
}
