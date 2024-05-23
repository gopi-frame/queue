package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/database"
	"github.com/gopi-frame/support/utils"
	"gorm.io/gorm"
)

// New new queue
func New(name string, connection *database.Connection, table string, failer queue.Failer) *Queue {
	q := new(Queue)
	q.name = name
	q.table = table
	q.connection = connection
	q.failer = failer
	q.fields = new(Job).Fields(table)
	return q
}

// Queue queue
type Queue struct {
	queue.AbstractQueue
	name       string
	table      string
	connection *database.Connection
	failer     queue.Failer

	fields JobDAO
}

// Name return queue name
func (q *Queue) Name() string {
	return q.name
}

// Count returns the count of pending jobs
func (q *Queue) Count() int64 {
	count, err := new(Job).UseConn(q.connection).
		UseTable(q.table).
		Where(q.fields.Queue.Eq(q.name)).
		Count()
	if err != nil {
		panic(err)
	}
	return count
}

// IsEmpty returns if the count of pending jobs is zero
func (q *Queue) IsEmpty() bool {
	_, err := new(Job).UseConn(q.connection).
		UseTable(q.table).
		Where(q.fields.Queue.Eq(q.name)).
		Where(q.fields.ReservedAt.Null()).
		Where(q.fields.AvaliableAt.Lte(time.Now())).
		First()
	return err == gorm.ErrRecordNotFound
}

// Enqueue pushes a job to queue
func (q *Queue) Enqueue(job queue.Job) bool {
	err := new(Job).UseConn(q.connection).
		UseTable(q.table).
		Create(&Job{
			JobUUID:     uuid.New(),
			Queue:       q.name,
			JobPayload:  job,
			AvaliableAt: utils.Ptr(time.Now().Add(job.Delay())),
		})
	if err != nil {
		panic(err)
	}
	return true
}

// Dequeue pops a job from queue
func (q *Queue) Dequeue() (queue.Job, bool) {
	model, err := new(Job).UseConn(q.connection).
		UseTable(q.table).
		Where(q.fields.Queue.Eq(q.name)).
		Where(q.fields.ReservedAt.Null()).
		Where(q.fields.AvaliableAt.Lte(time.Now())).
		Order(q.fields.AvaliableAt.Asc()).
		Take()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		return nil, false
	}
	return model.Payload(), true
}

// Remove removes a job from queue
func (q *Queue) Remove(job queue.Job) bool {
	_, err := new(Job).UseConn(q.connection).
		UseTable(q.table).
		Where(q.fields.Queue.Eq(q.name)).
		Where(q.fields.UUID.Eq(job.GetJob().UUID().String())).
		Delete()
	if err != nil {
		panic(err)
	}
	return true
}

// Ack acks a job
func (q *Queue) Ack(job queue.Job) {
	q.Remove(job)
}

// Fail handles a failed job
func (q *Queue) Fail(job queue.Job, err error) {
	if q.failer != nil {
		q.failer.Save(q.name, job, err)
	}
}

// Flush removes all failed jobs
func (q *Queue) Flush() {
	if q.failer != nil {
		q.failer.Flush(q.name)
	}
}

// Reload reloads all failed jobs into queue
func (q *Queue) Reload() {
	if q.failer != nil {
		q.failer.All(q.name).Each(func(key int, value queue.Job) bool {
			return q.Enqueue(value)
		})
	}
}

// Progress progress
func (q *Queue) Progress() (pending int64, executing int64) {
	err := q.connection.Transaction(func(tx *gorm.DB) error {
		var err error
		executing, err = new(Job).UseTx(tx).
			Where(q.fields.Queue.Eq(q.name)).
			Where(q.fields.ReservedAt.Null()).
			Count()
		pending, err = new(Job).UseTx(tx).
			Where(q.fields.Queue.Eq(q.name)).
			Where(q.fields.ReservedAt.Null()).
			Count()
		return err
	})
	if err != nil {
		panic(err)
	}
	return pending, executing
}
