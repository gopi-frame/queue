package empty

import (
	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/queue"
)

// Queue empty driver
type Queue struct{}

// Name return queue name
func (q *Queue) Name() string {
	return ""
}

// Count returns the count of pending jobs
func (q *Queue) Count() int64 {
	return 0
}

// IsEmpty returns if the count of pending jobs is zero
func (q *Queue) IsEmpty() bool {
	return true
}

// Enqueue pushes a job to queue
func (q *Queue) Enqueue(job queue.Job) bool {
	return true
}

// Dequeue pops a job from queue
func (q *Queue) Dequeue() (queue.Job, bool) {
	return nil, false
}

// Remove removes a job from queue
func (q *Queue) Remove(job queue.Job) bool {
	return true
}

// Ack acks a job
func (q *Queue) Ack(job queue.Job) {}

// Fail handles a failed job
func (q *Queue) Fail(job queue.Job, err error) {}

// Flush removes all failed jobs
func (q *Queue) Flush() {}

// Reload reloads all failed jobs into queue
func (q *Queue) Reload() {}

// Subscribe add a subscriber to queue events
func (q *Queue) Subscribe(subscriber queue.Subcriber) {}

// DispatchEvent dispatches specifia event
func (q *Queue) DispatchEvent(event event.Event) {}

// Progress progress
func (q *Queue) Progress() (pending int64, executing int64) {
	return 0, 0
}
