package memory

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

func NewMemoryJob(job queue.JobInterface, queue string) *MemoryJob {
	model := &MemoryJob{
		ID:      uuid.New(),
		Payload: job,
		Queue:   queue,
	}
	job.SetModel(model)
	return model
}

type MemoryJob struct {
	ID       uuid.UUID
	Queue    string
	Payload  queue.JobInterface
	Attempts int
}

func (d *MemoryJob) GetID() string {
	return d.ID.String()
}

func (d *MemoryJob) GetQueue() string {
	return d.Queue
}

func (d *MemoryJob) GetPayload() queue.JobInterface {
	return d.Payload
}

func (d *MemoryJob) GetAttempts() int {
	return d.Attempts
}
