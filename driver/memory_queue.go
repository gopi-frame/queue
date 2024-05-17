package driver

import (
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/event"
	"github.com/gopi-frame/queue/failed"
	"github.com/gopi-frame/queue/job"
	"github.com/gopi-frame/support/lists"
)

var _ queue.Queue = (*MemoryQueue)(nil)

// MemoryQueue memory workerpool driver
type MemoryQueue struct {
	Queue
	queue      string
	jobs       *lists.LinkedList[*job.MemoryJob]
	failedJobs queue.FailedJobProvider
}

// NewMemoryDriver creates a new memory driver
func NewMemoryDriver() *MemoryQueue {
	mq := new(MemoryQueue)
	mq.events = event.NewDispatcher()
	mq.jobs = lists.NewLinkedList[*job.MemoryJob]()
	mq.failedJobs = &failed.EmptyFailedJobProvider{}
	return mq
}

// Name queue
func (mq *MemoryQueue) Name() string {
	return mq.queue
}

// Count returns the count of pending jobs
func (mq *MemoryQueue) Count() (count int64) {
	mq.jobs.Lock()
	defer mq.jobs.Unlock()
	count = mq.jobs.Count()
	return
}

// IsEmpty returns if the count of pending jobs is zero
func (mq *MemoryQueue) IsEmpty() (isEmpty bool) {
	mq.jobs.Lock()
	defer mq.jobs.Unlock()
	isEmpty = mq.jobs.IsEmpty()
	return
}

// Enqueue pushes a job to queue
func (mq *MemoryQueue) Enqueue(item queue.Job) (success bool) {
	mq.jobs.Lock()
	defer mq.jobs.Unlock()
	mq.jobs.Push(job.NewMemoryJob(item))
	success = true
	return
}

// Dequeue pops a job from queue
func (mq *MemoryQueue) Dequeue() (result queue.Job, success bool) {
	mq.jobs.Lock()
	defer mq.jobs.Unlock()
	if mq.jobs.IsEmpty() {
		result, success = nil, false
		return
	}
	job, found := mq.jobs.FirstWhere(func(value *job.MemoryJob) bool {
		return !value.Executing()
	})
	result, success = job.Payload(), found
	return
}

// Remove removes a job from queue
func (mq *MemoryQueue) Remove(value queue.Job) (success bool) {
	mq.jobs.Lock()
	defer mq.jobs.Unlock()
	mq.jobs.RemoveWhere(func(item *job.MemoryJob) bool {
		return item.UUID() == value.GetJob().UUID()
	})
	success = true
	return
}

// Ack acks a job
func (mq *MemoryQueue) Ack(job queue.Job) {
	mq.Remove(job)
}

// Fail handles a failed job
func (mq *MemoryQueue) Fail(job queue.Job, err error) {
	mq.failedJobs.Save(mq.queue, job, err)
}

// Flush removes all failed jobs
func (mq *MemoryQueue) Flush() {
	mq.failedJobs.Flush()
}

// Reload reloads all failed jobs into queue
func (mq *MemoryQueue) Reload() {
	mq.failedJobs.All().Each(func(key int, value queue.Job) bool {
		return mq.Enqueue(value)
	})
}

// Progress progress
func (mq *MemoryQueue) Progress() (executing int64, pending int64) {
	mq.jobs.Lock()
	defer mq.jobs.Unlock()
	executing = int64(mq.jobs.Where(func(value *job.MemoryJob) bool {
		return value.Executing()
	}).Count())
	pending = mq.Count() - executing
	return
}
