package queue

import (
	"time"

	"github.com/gopi-frame/contract/queue"
)

// Job abstract job
type Job struct {
	queue.Job
	job queue.Queueable
}

// Load set implemention
func (j *Job) Load(impl queue.Job) {
	j.Job = impl
}

// SetJob set queue job
func (j *Job) SetJob(job queue.Queueable) {
	j.job = job
}

// GetJob get queue job
func (j *Job) GetJob() queue.Queueable {
	return j.job
}

// Delay get job delay
func (j *Job) Delay() time.Duration {
	return 0
}

// Timeout get job timeout
func (j *Job) Timeout() time.Duration {
	return 0
}

// MaxAttempts get job max attempts
func (j *Job) MaxAttempts() uint {
	return 1
}

// RetryDelay get job retry delay
func (j *Job) RetryDelay() time.Duration {
	return 0
}

// MaxRetryDelay get job max retry delay
func (j *Job) MaxRetryDelay() time.Duration {
	return 0
}

// RetryDelayStep get job retry delay step
func (j *Job) RetryDelayStep() time.Duration {
	return 0
}

// Failed failed
func (j *Job) Failed(err error) {}
