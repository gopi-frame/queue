package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/dao"
	"gorm.io/gen/field"
)

// Job job model
type Job struct {
	dao.DAO[Job] `gorm:"-"`

	JobID       uint64     `gorm:"column:id;autoIncrement;primaryKey;not null"`
	JobUUID     uuid.UUID  `gorm:"column:uuid;uniqueIndex;not null"`
	Queue       string     `gorm:"column:queue;not null"`
	JobPayload  queue.Job  `gorm:"column:payload;serializer:json;not null"`
	JobAttempts uint8      `gorm:"column:attempts;default:0"`
	ReservedAt  *time.Time `gorm:"column:reserved_at"`
	AvaliableAt *time.Time `gorm:"column:avaliable_at;not null"`
	CreatedAt   *time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

// JobDAO job dao
type JobDAO struct {
	dao.DAO[Job]

	ID          field.Uint64
	UUID        field.String
	Queue       field.String
	Payload     field.String
	Attempts    field.Uint8
	ReservedAt  field.Time
	AvaliableAt field.Time
	CreatedAt   field.Time
}

// ID id
func (j *Job) ID() uint64 {
	return j.JobID
}

// UUID uuid
func (j *Job) UUID() uuid.UUID {
	return j.JobUUID
}

// Attempts attempts
func (j *Job) Attempts() uint {
	return uint(j.JobAttempts)
}

// Payload payload
func (j *Job) Payload() queue.Job {
	return j.JobPayload
}

// Fields fields
func (j Job) Fields(table string) JobDAO {
	return JobDAO{
		DAO: *(dao.DAO[Job]{}).UseTable(table),

		ID:          field.NewUint64(table, "id"),
		UUID:        field.NewString(table, "uuid"),
		Queue:       field.NewString(table, "queue"),
		Payload:     field.NewString(table, "payload"),
		Attempts:    field.NewUint8(table, "attempts"),
		ReservedAt:  field.NewTime(table, "reserved_at"),
		AvaliableAt: field.NewTime(table, "avaliable_at"),
		CreatedAt:   field.NewTime(table, "created_at"),
	}
}
