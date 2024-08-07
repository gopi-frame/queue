package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

// NewJob creates a new database queue job
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

// Job is a database queue job, it is a wrapper around the [queue.Job]
type Job struct {
	ID          uuid.UUID           `gorm:"column:id;primaryKey" json:"id"`
	Queue       string              `gorm:"column:queue" json:"queue"`
	Payload     queue.Job           `gorm:"column:payload;serializer:json" json:"payload"`
	Attempts    int                 `gorm:"column:attempts" json:"attempts"`
	ReservedAt  sql.Null[time.Time] `gorm:"column:reserved_at" json:"reserved_at"`
	AvailableAt time.Time           `gorm:"column:available_at" json:"available_at"`
	CreatedAt   time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"column:updated_at" json:"updated_at"`
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

// GetAttempts returns how many times the job has been attempted
func (d *Job) GetAttempts() int {
	return d.Attempts
}
