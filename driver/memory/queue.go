package memory

import (
	"github.com/gopi-frame/collection/queue"
	queuecontract "github.com/gopi-frame/contract/queue"
	"time"
)

// Queue is a memory queue
type Queue struct {
	name string
	jobs *queue.PriorityQueue[*Job]
}

type Comparator struct{}

func (Comparator) Compare(a, b *Job) int {
	if a.AvailableAt.Before(b.AvailableAt) {
		return -1
	} else if a.AvailableAt.After(b.AvailableAt) {
		return 1
	}
	return 0
}

// NewQueue creates a new memory queue
func NewQueue(name string) *Queue {
	return &Queue{
		name: name,
		jobs: queue.NewPriorityQueue[*Job](Comparator{}),
	}
}

// Name returns the queue name
func (q *Queue) Name() string {
	return q.name
}

// Empty returns true if the queue is empty
func (q *Queue) Empty() bool {
	q.jobs.RLock()
	defer q.jobs.RUnlock()
	return q.jobs.IsEmpty()
}

// Count returns the number of jobs in the queue
func (q *Queue) Count() int64 {
	q.jobs.RLock()
	defer q.jobs.RUnlock()
	return q.jobs.Count()
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(job queuecontract.Job) (queuecontract.Job, bool) {
	q.jobs.Lock()
	defer q.jobs.Unlock()
	model := NewJob(job, q.name)
	return model.Payload, q.jobs.Enqueue(model)
}

// Dequeue removes a job from the queue
func (q *Queue) Dequeue() (queuecontract.Job, bool) {
	q.jobs.Lock()
	defer q.jobs.Unlock()
	model, ok := q.jobs.Peek()
	if ok {
		if model.AvailableAt.After(time.Now()) {
			return nil, false
		}
		model, _ = q.jobs.Dequeue()
		return model.Payload, true
	}
	return nil, false
}

// Remove removes a job from the queue
func (q *Queue) Remove(job queuecontract.Job) {
	if model := job.GetQueueable(); model != nil {
		q.jobs.Lock()
		defer q.jobs.Unlock()
		q.jobs.RemoveWhere(func(value *Job) bool {
			if value.ID.String() == model.GetID() {
				return true
			}
			return false
		})
	}
}

// Ack acknowledges a job
func (q *Queue) Ack(_ queuecontract.Job) {}

// Release releases a job and adds it back to the queue for the next attempt.
func (q *Queue) Release(job queuecontract.Job) {
	if model := job.GetQueueable(); model != nil {
		q.jobs.Lock()
		defer q.jobs.Unlock()
		model := model.(*Job)
		model.Attempts++
		model.AvailableAt = time.Now()
		q.jobs.Enqueue(model)
	}
}
