package memory

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

func NewJob(job queue.Job, queue string) *Job {
	model := &Job{
		ID:      uuid.New(),
		Payload: job,
		Queue:   queue,
	}
	job.SetModel(model)
	return model
}

type Job struct {
	ID       uuid.UUID
	Queue    string
	Payload  queue.Job
	Attempts int
}

func (d *Job) GetID() string {
	return d.ID.String()
}

func (d *Job) GetQueue() string {
	return d.Queue
}

func (d *Job) GetPayload() queue.Job {
	return d.Payload
}

func (d *Job) GetAttempts() int {
	return d.Attempts
}
