package memory

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/future"
	squeue "github.com/gopi-frame/support/queue"
)

func NewQueue(name string) *Queue {
	return &Queue{
		name: name,
		size: new(atomic.Int64),
		jobs: squeue.NewQueue[*Job](),
	}
}

type Queue struct {
	name string
	size *atomic.Int64
	jobs *squeue.Queue[*Job]
}

func (q *Queue) Empty() bool {
	return q.size.Load() == 0
}

func (q *Queue) Count() int64 {
	return q.size.Load()
}

func (q *Queue) Enqueue(job queue.JobInterface) {
	q.jobs.Lock()
	defer q.jobs.RUnlock()
	if q.jobs.Enqueue(NewJob(job, q.name)) {
		q.size.Add(1)
	}
}

func (q *Queue) Dequeue() queue.JobInterface {
	q.jobs.RLock()
	defer q.jobs.RUnlock()
	job, ok := q.jobs.Dequeue()
	if ok {
		q.size.Add(-1)
		return job.Payload
	}
	return nil
}

func (q *Queue) Remove(job queue.JobInterface) {
	if job.GetModel() == nil {
		return
	}
	q.jobs.Lock()
	defer q.jobs.Unlock()
	q.jobs.RemoveWhere(func(value *Job) bool {
		if value.ID.String() == job.GetModel().GetID() {
			q.size.Add(-1)
			return true
		}
		return false
	})
}

func (q *Queue) Ack(job queue.JobInterface) {}

func (q *Queue) Release(job queue.JobInterface, delay time.Duration) {
	if model := job.GetModel(); model == nil {
		return
	}
	model := new(Job)
	model.ID = uuid.New()
	model.Queue = q.name
	model.Payload = job
	model.Attempts = job.GetModel().GetAttempts()
	job.SetModel(model)
	future.Delay(func() any {
		q.jobs.Lock()
		defer q.jobs.Unlock()
		if q.jobs.Enqueue(model) {
			q.size.Add(1)
		}
		return nil
	}, delay)
}
