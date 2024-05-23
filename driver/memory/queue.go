package memory

import (
	"sync"

	"github.com/google/uuid"
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/support/queue"
)

// Queue driver
type Queue struct {
	sync.Mutex
	queuecontract.AbstractQueue
	name           string
	pendingItems   *queue.LinkedQueue[*Job]
	executingItems *lists.List[*Job]
}

// Name return queue name
func (q *Queue) Name() string {
	return q.name
}

// Count returns the count of pending jobs
func (q *Queue) Count() int64 {
	if q.TryLock() {
		defer q.Unlock()
	}
	return q.executingItems.Count() + q.pendingItems.Count()
}

// IsEmpty returns if the count of pending jobs is zero
func (q *Queue) IsEmpty() bool {
	return q.Count() == 0
}

// Enqueue pushes a job to queue
func (q *Queue) Enqueue(job queuecontract.Job) bool {
	if q.TryLock() {
		defer q.Unlock()
	}
	queueable := &Job{
		uuid:    uuid.New(),
		payload: job,
	}
	job.SetJob(queueable)
	return q.pendingItems.Enqueue(queueable)
}

// Dequeue pops a job from queue
func (q *Queue) Dequeue() (queuecontract.Job, bool) {
	if q.TryLock() {
		defer q.Unlock()
	}
	job, ok := q.pendingItems.Dequeue()
	if !ok {
		return nil, false
	}
	q.executingItems.Push(job)
	return job.Payload(), true
}

// Remove removes a job from queue
func (q *Queue) Remove(job queuecontract.Job) bool {
	if q.TryLock() {
		defer q.Unlock()
	}
	q.pendingItems.RemoveWhere(func(value *Job) bool {
		return value.UUID() == job.GetJob().UUID()
	})
	return true
}

// Ack acks a job
func (q *Queue) Ack(job queuecontract.Job) {
	q.executingItems.RemoveWhere(func(value *Job) bool {
		return value.UUID() == job.GetJob().UUID()
	})
}

// Fail handles a failed job
func (q *Queue) Fail(job queuecontract.Job, err error) {
	panic("not implemented") // TODO: Implement
}

// Flush removes all failed jobs
func (q *Queue) Flush() {
	panic("not implemented") // TODO: Implement
}

// Reload reloads all failed jobs into queue
func (q *Queue) Reload() {
	panic("not implemented") // TODO: Implement
}

// Progress progress
func (q *Queue) Progress() (pending int64, executing int64) {
	return q.Count(), 0
}
