package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"reflect"
	"sync"
	"time"

	"github.com/gopi-frame/contract/queue"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// column names for Job Table
const (
	ColumnID          = "id"
	ColumnQueue       = "queue"
	ColumnPayload     = "payload"
	ColumnReservedAt  = "reserved_at"
	ColumnAvailableAt = "available_at"
	ColumnCreatedAt   = "created_at"
	ColumnUpdatedAt   = "updated_at"
)

// Queue is database backed queue implementation
type Queue struct {
	mu sync.RWMutex
	db *gorm.DB

	name  string
	table string
	job   reflect.Type
}

// NewQueue returns a new queue.
func NewQueue(cfg *Config, opts ...Option) *Queue {
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			panic(err)
		}
	}
	if cfg.Table == "" {
		cfg.Table = DefaultJobTable
	}
	if err := cfg.Valid(); err != nil {
		panic(err)
	}
	jobType := reflect.Indirect(reflect.ValueOf(cfg.Job)).Type()
	return &Queue{
		db:    cfg.DB,
		name:  cfg.Name,
		table: cfg.Table,
		job:   jobType,
	}
}

// Name returns the queue name.
func (q *Queue) Name() string {
	return q.name
}

// Empty returns true if queue is empty.
func (q *Queue) Empty() bool {
	var model = make(map[string]any)
	result := q.db.Table(q.table).Where(clause.Eq{
		Column: clause.Column{Name: ColumnQueue},
		Value:  q.name,
	}).Where(clause.Eq{
		Column: clause.Column{Name: ColumnReservedAt},
	}).Limit(1).Take(&model)
	if err := result.Error; err == nil {
		return false
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	} else {
		panic(err)
	}
}

// Count returns the number of jobs in the queue.
func (q *Queue) Count() int64 {
	var count int64
	result := q.db.Model(new(Job)).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Table(q.table)
	}).Where(clause.Eq{
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

// Enqueue adds a job to the queue.
func (q *Queue) Enqueue(job queue.Job) (queue.Job, bool) {
	model := NewJob(q.name, job)
	result := q.db.Table(q.table).Create(model)
	if err := result.Error; err != nil {
		panic(err)
	}
	return model.Payload, true
}

// Dequeue removes a job from the queue and returns it.
func (q *Queue) Dequeue() (queue.Job, bool) {
	q.mu.RLock()
	defer q.mu.RLock()

	var dest struct {
		ID          uuid.UUID
		Queue       string
		Payload     string
		ReservedAt  sql.Null[time.Time]
		Attempts    int
		AvailableAt time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	var model = new(Job)

	err := q.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Table(q.table).Where(clause.Eq{
			Column: clause.Column{Name: ColumnQueue},
			Value:  q.name,
		}).Where(clause.Eq{
			Column: clause.Column{Name: ColumnReservedAt},
		}).Where(clause.Lte{
			Column: clause.Column{Name: ColumnAvailableAt},
			Value:  time.Now(),
		}).Order(clause.OrderByColumn{
			Column: clause.Column{Name: ColumnCreatedAt},
		}).Take(&dest)
		if err := result.Error; err != nil {
			return err
		}

		model.ID = dest.ID
		model.Queue = dest.Queue
		model.Attempts = dest.Attempts
		model.AvailableAt = dest.AvailableAt
		model.CreatedAt = dest.CreatedAt
		model.UpdatedAt = dest.UpdatedAt
		model.ReservedAt = sql.Null[time.Time]{V: time.Now(), Valid: true}

		var payload = reflect.New(q.job).Interface()
		if err := json.Unmarshal([]byte(dest.Payload), payload); err != nil {
			return err
		}
		model.Payload = payload.(queue.Job)
		model.Payload.SetQueueable(model)

		return tx.Save(model).Error
	})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
		return nil, false
	}
	model.Payload.SetQueueable(model)
	return model.Payload, true
}

// Remove removes a job from the queue.
func (q *Queue) Remove(job queue.Job) {
	if model := job.GetQueueable(); model != nil {
		result := q.db.Table(q.table).Where(clause.Eq{
			Column: ColumnID,
			Value:  model.GetID(),
		}).Delete(nil)
		if err := result.Error; err != nil {
			panic(err)
		}
	}
}

// Ack acknowledges a job.
func (q *Queue) Ack(job queue.Job) {
	q.Remove(job)
}

// Release releases a job back to the queue.
func (q *Queue) Release(job queue.Job) {
	if model := job.GetQueueable(); model != nil {
		err := q.db.Transaction(func(tx *gorm.DB) error {
			model := model.(*Job)
			model.Attempts++
			model.AvailableAt = time.Now()
			return tx.Save(model).Error
		})
		if err != nil {
			panic(err)
		}
	}
}
