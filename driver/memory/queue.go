package memory

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/future"
	squeue "github.com/gopi-frame/support/queue"
)

type Queue struct {
	name string
	size atomic.Int64
	jobs *squeue.Queue[*MemoryJob]
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
	q.jobs.Enqueue(NewMemoryJob(job, q.name))
}

func (q *Queue) Dequeue() queue.JobInterface {
	q.jobs.RLock()
	defer q.jobs.RUnlock()
	job, ok := q.jobs.Dequeue()
	if ok {
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
	q.jobs.RemoveWhere(func(value *MemoryJob) bool {
		return value.ID.String() == job.GetModel().GetID()
	})
}

func (q *Queue) Ack(job queue.JobInterface) {}

func (q *Queue) Release(job queue.JobInterface, delay time.Duration) {
	if model := job.GetModel(); model == nil {
		return
	}
	model := new(MemoryJob)
	model.ID = uuid.New()
	model.Queue = q.name
	model.Payload = job
	model.Attempts = job.GetModel().GetAttempts()
	job.SetModel(model)
	future.Delay(func() any {
		q.jobs.Lock()
		defer q.jobs.Unlock()
		q.jobs.Enqueue(model)
		return nil
	}, delay)
}
