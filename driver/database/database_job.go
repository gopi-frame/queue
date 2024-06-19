package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

const (
	DefaultJobTable = "jobs"
)

const (
	ColumnID          = "id"
	ColumnQueue       = "queue"
	ColumnPayload     = "payload"
	ColumnReservedAt  = "reserved_at"
	ColumnAvaliableAt = "avaliable_at"
	ColumnCreatedAt   = "created_at"
	ColumnUpdatedAt   = "updated_at"
)

func NewDatabaseJob(queue string, payload queue.JobInterface) *DatabaseJob {
	queueable := &DatabaseJob{
		ID:          uuid.New(),
		Queue:       queue,
		Payload:     payload,
		AvaliableAt: time.Now().Add(payload.GetDelay()),
	}
	payload.SetModel(queueable)
	return queueable
}

type DatabaseJob struct {
	ID          uuid.UUID           `gorm:"column:id;primaryKey" json:"id"`
	Queue       string              `gorm:"column:queue" json:"queue"`
	Payload     queue.JobInterface  `gorm:"column:payload;serializer:json" json:"payload"`
	Attempts    int                 `gorm:"column:attempts" json:"attempts"`
	ReservedAt  sql.Null[time.Time] `gorm:"column:reserved_at" json:"reserved_at"`
	AvaliableAt time.Time           `gorm:"column:avaliable_at" json:"avaliable_at"`
	CreatedAt   time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"column:updated_at" json:"updated_at"`
}

func (d *DatabaseJob) GetID() string {
	return d.ID.String()
}

func (d *DatabaseJob) GetQueue() string {
	return d.Queue
}

func (d *DatabaseJob) GetPayload() queue.JobInterface {
	return d.Payload
}

func (d *DatabaseJob) GetAttempts() int {
	return d.Attempts
}
