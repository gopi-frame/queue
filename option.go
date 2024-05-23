package queue

import (
	"time"

	"github.com/gopi-frame/contract/queue"
)

// Default option values
const (
	DefaultMaxWorkerNum         = 10
	DefaultWorkerCreationBatch  = 1
	DefaultWorkerMaxIdleTime    = 5 * time.Minute
	DefaultWorkerMaxStoppedTime = 10 * time.Minute
)

// Option worker pool option fn
type Option func(*Queue)

var noneOption = func(*Queue) {}

// WorkerNum sets count of workers, default is 10
func WorkerNum(count int) Option {
	return func(wp *Queue) {
		if count <= 0 {
			count = DefaultMaxWorkerNum
		}
		wp.workerNum = count
	}
}

// WorkerBatch sets count of worker creation batches, default is 1
func WorkerBatch(batch int) Option {
	return func(wp *Queue) {
		if batch <= 0 {
			batch = DefaultWorkerCreationBatch
		}
		wp.workerBatch = batch
	}
}

// WorkerMaxIdleTime sets max idle time of a worker, default is 5min
func WorkerMaxIdleTime(d time.Duration) Option {
	return func(wp *Queue) {
		if d < 0 {
			d = DefaultWorkerMaxIdleTime
		}
		wp.workerMaxIdleTime = d
	}
}

// WorkerMaxStoppedTime sets max stopped time of a worker, default is 10min
func WorkerMaxStoppedTime(d time.Duration) Option {
	return func(wp *Queue) {
		if d < 0 {
			d = DefaultWorkerMaxStoppedTime
		}
		wp.workerMaxStoppedTime = d
	}
}

// Subscriber adds a subscriber to queue events
func Subscriber(subscriber queue.Subcriber) Option {
	if subscriber == nil {
		return noneOption
	}
	return func(wp *Queue) {
		wp.Subscribe(subscriber)
	}
}
