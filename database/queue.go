package database

import (
	"database/sql"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/queue/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewQueue(cfg *Config) *Queue {
	return &Queue{
		db:    cfg.db,
		name:  cfg.name,
		table: cfg.table,
	}
}

type Queue struct {
	mu sync.RWMutex
	db *gorm.DB

	name  string
	table string
}

func (q *Queue) Empty() bool {
	model := new(Job)
	result := q.db.Where(clause.Eq{
		Column: clause.Column{Name: ColumnQueue},
		Value:  q.name,
	}).Where(clause.Eq{
		Column: clause.Column{Name: ColumnReservedAt},
	}).Take(model)
	if err := result.Error; err == nil {
		return false
	} else if err == gorm.ErrRecordNotFound {
		return true
	} else {
		panic(err)
	}
}

func (q *Queue) Count() int64 {
	var count int64
	result := q.db.Where(clause.Eq{
		Column: clause.Column{Name: ColumnQueue},
		Value:  q.name,
	}).Where(clause.Eq{
		Column: clause.Column{Name: ColumnReservedAt},
	}).Count(&count)
	if err := result.Error; err != nil {
		panic(err)
	}
	return count
}

func (q *Queue) Enqueue(job driver.Job) {
	model := NewJob(q.name, job)
	result := q.db.Table(q.table).Create(model)
	if err := result.Error; err != nil {
		panic(err)
	}
}

func (q *Queue) Dequeue() driver.Job {
	q.mu.RLock()
	defer q.mu.RLock()
	model := new(Job)
	err := q.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Table(q.table).Where(clause.Eq{
			Column: clause.Column{Name: ColumnQueue},
			Value:  q.name,
		}).Where(clause.Eq{
			Column: clause.Column{Name: ColumnReservedAt},
		}).Where(clause.Lte{
			Column: clause.Column{Name: ColumnAvaliableAt},
			Value:  time.Now(),
		}).Order(clause.OrderByColumn{
			Column: clause.Column{Name: ColumnCreatedAt},
		}).Take(model)
		if err := result.Error; err != nil {
			return err
		}
		model.Attempts++
		model.ReservedAt = sql.Null[time.Time]{V: time.Now(), Valid: true}
		return tx.Save(model).Error
	})
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return nil
	}
	return model.Payload
}

func (q *Queue) Remove(job driver.Job) {
	if job.GetModel() == nil {
		return
	}
	result := q.db.Where(clause.Eq{
		Column: ColumnID,
		Value:  job.GetModel().GetID(),
	}).Delete(nil)
	if err := result.Error; err != nil {
		panic(err)
	}
}

func (q *Queue) Ack(job driver.Job) {
	q.Remove(job)
}

func (q *Queue) Release(job driver.Job, delay time.Duration) {
	err := q.db.Transaction(func(tx *gorm.DB) error {
		if model := job.GetModel(); model != nil {
			if err := tx.Delete(model).Error; err != nil {
				return err
			}
		}
		model := new(Job)
		model.ID = uuid.New()
		model.Queue = q.name
		model.Payload = job
		model.Attempts = job.GetModel().GetAttempts()
		model.AvaliableAt = time.Now().Add(delay)
		job.SetModel(model)
		return tx.Create(model).Error
	})
	if err != nil {
		panic(err)
	}
}
