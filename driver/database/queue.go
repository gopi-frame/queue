package database

import (
	"time"

	"github.com/gopi-frame/contract/queue"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultJobTable = "jobs"

func NewDatabaseQueue(db *gorm.DB, table string, queue string) *DatabaseQueue {
	return &DatabaseQueue{
		queue: queue,
		table: table,
		db:    db,
	}
}

type DatabaseQueue struct {
	queue string
	table string
	db    *gorm.DB
}

func (q *DatabaseQueue) Queue() string {
	return q.queue
}

func (q *DatabaseQueue) Count() (int64, error) {
	var count int64
	if err := q.db.Table(q.table).Where(clause.Eq{
		Column: clause.Column{Name: "queue"},
		Value:  q.queue,
	}).Where(clause.Lte{
		Column: clause.Column{Name: "avaliable_at"},
		Value:  time.Now(),
	}).Where(clause.Eq{Column: clause.Column{Name: "reserved_at"}}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *DatabaseQueue) Empty() (bool, error) {
	var model = new(DatabaseJob)
	if err := q.db.Table(q.table).Where(clause.Eq{
		Column: clause.Column{Name: "queue"},
		Value:  q.queue,
	}).Where(clause.Lte{
		Column: clause.Column{Name: "avaliable_at"},
		Value:  time.Now(),
	}).Where(clause.Eq{Column: clause.Column{Name: "reserved_at"}}).Take(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (q *DatabaseQueue) Enqueue(job queue.Job) error {
	return q.db.Table(q.table).Where(clause.Eq{
		Column: clause.Column{Name: "queue"},
		Value:  q.queue,
	}).Create(NewDatabaseJob(job, q.queue)).Error
}

func (q *DatabaseQueue) Dequeue() (queue.Job, error) {
	var model = new(DatabaseJob)
	if err := q.db.Table(q.table).Where(clause.Eq{
		Column: clause.Column{Name: "queue"},
		Value:  q.queue,
	}).Where(clause.Lte{
		Column: clause.Column{Name: "avaliable_at"},
		Value:  time.Now(),
	}).Where(clause.Eq{Column: clause.Column{Name: "reserved_at"}}).First(model).Error; err != nil {
		return nil, err
	}
	return model.GetPayload(), nil
}

func (q *DatabaseQueue) Remove(job queue.Job) error {
	queueable := job.Queueable()
	if queueable == nil {
		return nil
	}
	return q.db.Table(q.table).Where(clause.Eq{
		Column: clause.Column{Name: "queue"},
		Value:  q.queue,
	}).Where(clause.Eq{
		Column: clause.Column{Name: "uuid"},
		Value:  queueable.GetUUID().String(),
	}).Delete(nil).Error
}
