package database

import (
	"time"

	"github.com/gopi-frame/contract/exception"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/dao"
	"gorm.io/gen/field"
)

// FailedJob failed job model
type FailedJob struct {
	ID        uint64              `gorm:"column:id;autoIncrement;primaryKey;not null"`
	Queue     string              `gorm:"column:queue;not null"`
	Payload   queue.Queueable     `gorm:"column:payload;serializer:json;not null"`
	Exception exception.Throwable `gorm:"column:exception"`
	FailedAt  time.Time           `gorm:"column:failed_at;not null;autoCreateTime"`
}

// FailedJobDAO failed job dao
type FailedJobDAO struct {
	dao.DAO[FailedJob]

	ID        field.Uint64
	Queue     field.String
	Payload   field.String
	Exception field.String
	FailedAt  field.Time
}

// Query query
func (f *FailedJob) Query(table string) FailedJobDAO {
	return FailedJobDAO{
		DAO: *dao.New[FailedJob]().UseTable(table),

		ID:        field.NewUint64(table, "id"),
		Queue:     field.NewString(table, "queue"),
		Payload:   field.NewString(table, "payload"),
		Exception: field.NewString(table, "exception"),
		FailedAt:  field.NewTime(table, "failed_at"),
	}
}
