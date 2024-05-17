package event

// NewProgressUpdated new progress updated
func NewProgressUpdated(pending, executing int64) *ProgressUpdated {
	return &ProgressUpdated{
		pending,
		executing,
	}
}

// ProgressUpdated progress updated event
type ProgressUpdated struct {
	Pending   int64
	Executing int64
}

// Topic topic
func (event ProgressUpdated) Topic() string {
	return ProgressUpdatedTopic
}
