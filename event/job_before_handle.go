package event

// JobBeforeHandle is a job before handle event
type JobBeforeHandle struct {
	Name string
	ID   string
}

// NewJobBeforeHandle creates a new job before handle event
func NewJobBeforeHandle(ID string, name string) *JobBeforeHandle {
	return &JobBeforeHandle{ID: ID, Name: name}
}

// Topic returns the topic
func (JobBeforeHandle) Topic() string {
	return TopicJobBeforeHandle
}
