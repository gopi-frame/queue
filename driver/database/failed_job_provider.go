package database

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"gorm.io/gorm"
)

const DefaultFailedJobTable = "failed_jobs"

func NewDatabaseFailedJobProvider(db *gorm.DB, table string) *DatabaseFailedJobProvider {
	return &DatabaseFailedJobProvider{
		db:    db,
		table: table,
	}
}

type DatabaseFailedJobProvider struct {
	db    *gorm.DB
	table string
}

func (d *DatabaseFailedJobProvider) Find(uuid uuid.UUID) queue.FailedJob {
	query := new(DatabaseFailedJob).Query(d.db, d.table)
	model, err := query.Where(query.UUID.Eq(uuid.String())).Take()
	if err == nil {
		return model
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	panic(err)
}

func (d *DatabaseFailedJobProvider) All(queue string) (items []queue.FailedJob) {
	query := new(DatabaseFailedJob).Query(d.db, d.table)
	models, err := query.Where(query.Queue.Eq(queue)).Find()
	if err != nil {
		panic(err)
	}
	for _, model := range models {
		items = append(items, model)
	}
	return items
}

func (d *DatabaseFailedJobProvider) Forget(uuid uuid.UUID) {
	query := new(DatabaseFailedJob).Query(d.db, d.table)
	_, err := query.Where(query.UUID.Eq(uuid.String())).Delete()
	if err != nil {
		panic(err)
	}
}

func (d *DatabaseFailedJobProvider) Flush(queue string) {
	query := new(DatabaseFailedJob).Query(d.db, d.table)
	_, err := query.Where(query.Queue.Eq(queue)).Delete()
	if err != nil {
		panic(err)
	}
}
