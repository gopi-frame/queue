package database

import (
	"time"

	"github.com/google/uuid"
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/contract/support"
	"github.com/gopi-frame/database"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/support/lists"
	"gorm.io/datatypes"
	"gorm.io/gen"
)

// Failer failer
type Failer struct {
	table      string
	connection *database.Connection
}

// Save save
func (f *Failer) Save(queue string, job queuecontract.Job, err error) {
	e := new(FailedJob).Query(f.table).Create(&FailedJob{
		Queue:     queue,
		Payload:   job.GetJob(),
		Exception: exception.NewException(err.Error()),
		FailedAt:  time.Time{},
	})
	if e != nil {
		panic(e)
	}
}

// All all
func (f *Failer) All(queue string) support.List[queuecontract.Job] {
	query := new(FailedJob).Query(f.table)
	models, err := query.Where(query.Queue.Eq(queue)).Find()
	if err != nil {
		panic(err)
	}
	items := lists.NewList[queuecontract.Job]()
	for _, model := range models {
		items.Push(model.Payload.Payload())
	}
	return items
}

// Find find
func (f *Failer) Find(queue string, id uuid.UUID) queuecontract.Job {
	query := new(FailedJob).Query(f.table)
	model, err := query.Where(query.Queue.Eq(queue)).
		Where(gen.Cond(datatypes.JSONQuery("payload").Equals(id, "uuid"))...).
		First()
	if err != nil {
		panic(err)
	}
	return model.Payload.Payload()
}

// Forget forget
func (f *Failer) Forget(queue string, id uuid.UUID) {
	query := new(FailedJob).Query(f.table)
	_, err := query.Where(query.Queue.Eq(queue)).
		Where(gen.Cond(datatypes.JSONQuery("payload").Equals(id, "uuid"))...).
		Delete()
	if err != nil {
		panic(err)
	}
}

// Flush flush
func (f *Failer) Flush(queue string) {
	query := new(FailedJob).Query(f.table)
	_, err := query.Where(query.Queue.Eq(queue)).Delete()
	if err != nil {
		panic(err)
	}
}
