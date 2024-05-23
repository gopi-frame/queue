package memory

import (
	"github.com/google/uuid"
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/contract/support"
	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/support/maps"
)

// Failer failer
type Failer struct {
	items maps.Map[string, *lists.List[*FailedJob]]
}

// Save save
func (f *Failer) Save(queue string, job queuecontract.Job, err error) {
	if f.items.ContainsKey(queue) {
		items, _ := f.items.Get(queue)
		items.Push(&FailedJob{
			Payload:   job.GetJob(),
			Queue:     queue,
			Exception: err,
		})
	}
}

// All all
func (f *Failer) All(queue string) support.List[queuecontract.Job] {
	items := lists.NewList[queuecontract.Job]()
	jobs, ok := f.items.Get(queue)
	if !ok {
		return items
	}
	jobs.Each(func(index int, value *FailedJob) bool {
		items.Push(value.Payload.Payload())
		return true
	})
	return items
}

// Find find
func (f *Failer) Find(queue string, id uuid.UUID) queuecontract.Job {
	items, ok := f.items.Get(queue)
	if !ok {
		return nil
	}
	item := items.FirstWhereOr(func(item *FailedJob) bool {
		return item.Payload.UUID() == id
	}, nil)
	if item == nil {
		return nil
	}
	return item.Payload.Payload()
}

// Forget forget
func (f *Failer) Forget(queue string, id uuid.UUID) {
	items, ok := f.items.Get(queue)
	if !ok {
		return
	}
	items.RemoveWhere(func(item *FailedJob) bool {
		return item.Payload.UUID() == id
	})
}

// Flush flush
func (f *Failer) Flush(queue string) {
	f.items.Remove(queue)
}
