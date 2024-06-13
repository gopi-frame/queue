package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
)

func NewDatabaseJob(job queue.Job, queue string) *DatabaseJob {
	queueable := &DatabaseJob{
		UUID:        uuid.New(),
		Queue:       queue,
		Payload:     job,
		AvaliableAt: job.AvalidableAt(),
	}
	job.SetQueueable(queueable)
	return queueable
}

type DatabaseJob struct {
	ID          uint64    `gorm:"column:id;autoIncrement;primaryKey;not null"`
	UUID        uuid.UUID `gorm:"column:uuid;uniqueKey"`
	Queue       string    `gorm:"column:queue;index"`
	Payload     queue.Job `gorm:"column:payload;serializer:json;not null"`
	Attempts    uint      `gorm:"column:attempts;not null"`
	AvaliableAt time.Time `gorm:"column:avaliable_at;not null"`
	ReservedAt  time.Time `gorm:"column:reserved_at"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (job *DatabaseJob) GetID() uint64         { return job.ID }
func (job *DatabaseJob) GetUUID() uuid.UUID    { return job.UUID }
func (job *DatabaseJob) GetQueue() string      { return job.Queue }
func (job *DatabaseJob) GetPayload() queue.Job { return job.Payload }
func (job *DatabaseJob) GetAttempts() uint     { return job.Attempts }

func (job *DatabaseJob) Fire() {
	exception.Try(func() {
		if err := job.GetPayload().Handle(); err != nil {
			panic(err)
		}
	}).CatchAll(func(err error) {
		job.Fail(err)
	})
}

func (job *DatabaseJob) Fail(err error) {
	job.GetPayload().Failed(err)
}
