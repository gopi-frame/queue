// Package queue is a queue package for gopi-frame.
package queue

import (
	"github.com/gopi-frame/contract/eventbus"
	"github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue/event"
	"time"
)

// Queue is a wrapper around [queue.Queue]
// It consumes jobs from queue and dispatches them to workers
type Queue struct {
	queue.Queue
	workernum int
	workers   chan queue.Worker
	stop      chan struct{}
	eventbus  eventbus.Bus
	logger    logger.Logger

	startedAt time.Time
}

// NewQueue creates a new queue
func NewQueue(q queue.Queue, workernum int) *Queue {
	return &Queue{
		Queue:     q,
		workernum: workernum,
		workers:   make(chan queue.Worker, workernum),
		stop:      make(chan struct{}),
	}
}

// Run starts the queue
func (q *Queue) Run() {
	defer func() {
		q.fire(event.NewQueueAfterStop(q.Name(), q.Uptime()))
	}()
	q.fire(event.NewQueueBeforeRun(q.Name()))
	q.startedAt = time.Now()
	for i := 0; i < q.workernum; i++ {
		worker := NewWorker(q.workers, q.Queue, q.eventbus)
		q.workers <- worker
	}
	for {
		select {
		case <-q.stop:
			q.workers = make(chan queue.Worker, q.workernum)
			return
		default:
			job, ok := q.Dequeue()
			if ok {
				// pop an idle worker, or block until one is available
				worker := <-q.workers
				go worker.Handle(job)
			}
		}
	}
}

// Stop stops the queue
func (q *Queue) Stop() {
	q.fire(event.NewQueueBeforeStop(q.Name()))
	close(q.stop)
}

// Uptime returns the uptime of the queue
func (q *Queue) Uptime() time.Duration {
	if q.startedAt.IsZero() {
		return 0
	}
	return time.Since(q.startedAt)
}

func (q *Queue) fire(e eventbus.Event) {
	if q.eventbus != nil {
		err := q.eventbus.Dispatch(e)
		if err != nil {
			return
		}
	}
}
