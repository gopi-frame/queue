package event

// JobAfterRelease is a job after release event
type JobAfterRelease struct {
	// Name is queue name
	Name string
	// ID is job ID
	ID string
	// Cause is the error that caused the job to be released
	Cause error
	// Attempts is how many times the job has been attempted
	Attempts int
}

// NewJobAfterRelease creates a new job after release event
func NewJobAfterRelease(name string, cause error, id string, attempts int) *JobAfterRelease {
	return &JobAfterRelease{Name: name, Cause: cause, ID: id, Attempts: attempts}
}

func (JobAfterRelease) Topic() string {
	return TopicJobAfterRelease
}
