package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/dao"
	"gorm.io/gen/field"
)

func NewDatabaseFailedJob(job queue.Queueable, exception error) *DatabaseFailedJob {
	return &DatabaseFailedJob{
		UUID:      uuid.New(),
		Queue:     job.GetQueue(),
		Payload:   job.GetPayload(),
		Exception: exception.Error(),
	}
}

type DatabaseFailedJobDAO struct {
	dao.DAO[DatabaseFailedJob]

	ID        field.Uint64
	UUID      field.UUID
	Queue     field.String
	Payload   dao.JSON
	Exception field.String
	FailedAt  field.Time
}

type DatabaseFailedJob struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;not null"`
	UUID      uuid.UUID `gorm:"column:uuid;primaryKey"`
	Queue     string    `gorm:"column:queue"`
	Payload   queue.Job `gorm:"column:payload;serializer:json"`
	Exception string    `gorm:"column:exception"`
	FailedAt  time.Time `gorm:"column:failed_at;autoCreateTime"`
}

func (job *DatabaseFailedJob) GetID() uint64         { return job.ID }
func (job *DatabaseFailedJob) GetUUID() uuid.UUID    { return job.UUID }
func (job *DatabaseFailedJob) GetQueue() string      { return job.Queue }
func (job *DatabaseFailedJob) GetPayload() queue.Job { return job.Payload }
func (job *DatabaseFailedJob) GetException() string  { return job.Exception }
