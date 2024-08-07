package memory

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"time"
)

// Job is a memory queue job
type Job struct {
	ID          uuid.UUID
	Queue       string
	Payload     queue.Job
	Attempts    int
	AvailableAt time.Time
}

// NewJob creates a new memory queue job
func NewJob(job queue.Job, queue string) *Job {
	model := &Job{
		ID:          uuid.New(),
		Payload:     job,
		Queue:       queue,
		AvailableAt: time.Now().Add(job.GetDelay()),
	}
	job.SetQueueable(model)
	return model
}

// GetID returns the job ID
func (d *Job) GetID() string {
	return d.ID.String()
}

// GetQueue returns the job queue
func (d *Job) GetQueue() string {
	return d.Queue
}

// GetPayload returns the job payload
func (d *Job) GetPayload() queue.Job {
	return d.Payload
}

// GetAttempts returns the job attempts
func (d *Job) GetAttempts() int {
	return d.Attempts
}
