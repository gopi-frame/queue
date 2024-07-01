package memory

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/queue/driver"
)

func NewJob(job driver.Job, queue string) *Job {
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
	Payload  driver.Job
	Attempts int
}

func (d *Job) GetID() string {
	return d.ID.String()
}

func (d *Job) GetQueue() string {
	return d.Queue
}

func (d *Job) GetPayload() driver.Job {
	return d.Payload
}

func (d *Job) GetAttempts() int {
	return d.Attempts
}
