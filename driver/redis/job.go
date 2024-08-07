package redis

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"time"
)

// Job is redis queueable job wrapper
type Job struct {
	ID          uuid.UUID `json:"id"`
	Queue       string    `json:"queue"`
	Payload     queue.Job `json:"payload"`
	Attempts    int       `json:"attempts"`
	AvailableAt time.Time `json:"available_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewJob creates new redis queueable job
func NewJob(queue string, payload queue.Job) *Job {
	queueable := &Job{
		ID:          uuid.New(),
		Queue:       queue,
		Payload:     payload,
		AvailableAt: time.Now().Add(payload.GetDelay()),
	}
	payload.SetQueueable(queueable)
	return queueable
}

// GetID returns the job ID
func (j *Job) GetID() string {
	return j.ID.String()
}

// GetQueue returns the job queue
func (j *Job) GetQueue() string {
	return j.Queue
}

// GetPayload returns the job payload
func (j *Job) GetPayload() queue.Job {
	return j.Payload
}

// GetAttempts returns how many times the job has been attempted
func (j *Job) GetAttempts() int {
	return j.Attempts
}
