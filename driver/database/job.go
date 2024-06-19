package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

func NewJob(queue string, payload queue.JobInterface) *Job {
	queueable := &Job{
		ID:          uuid.New(),
		Queue:       queue,
		Payload:     payload,
		AvaliableAt: time.Now().Add(payload.GetDelay()),
	}
	payload.SetModel(queueable)
	return queueable
}

type Job struct {
	ID          uuid.UUID           `gorm:"column:id;primaryKey" json:"id"`
	Queue       string              `gorm:"column:queue" json:"queue"`
	Payload     queue.JobInterface  `gorm:"column:payload;serializer:json" json:"payload"`
	Attempts    int                 `gorm:"column:attempts" json:"attempts"`
	ReservedAt  sql.Null[time.Time] `gorm:"column:reserved_at" json:"reserved_at"`
	AvaliableAt time.Time           `gorm:"column:avaliable_at" json:"avaliable_at"`
	CreatedAt   time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"column:updated_at" json:"updated_at"`
}

func (d *Job) GetID() string {
	return d.ID.String()
}

func (d *Job) GetQueue() string {
	return d.Queue
}

func (d *Job) GetPayload() queue.JobInterface {
	return d.Payload
}

func (d *Job) GetAttempts() int {
	return d.Attempts
}
