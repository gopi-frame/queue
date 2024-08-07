// Package queue is a queue package for gopi-frame.
package queue

import (
	"github.com/gopi-frame/contract/eventbus"
	"github.com/gopi-frame/contract/logger"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue/event"
	"time"
)

// DefaultWorkerNum is the default number of workers
const DefaultWorkerNum = 3

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
func NewQueue(buffer queue.Queue, opts ...Option) *Queue {
	q := &Queue{
		Queue: buffer,
		stop:  make(chan struct{}),
	}
	for _, opt := range opts {
		if err := opt(q); err != nil {
			panic(err)
		}
	}
	return q
}

type Option func(q *Queue) error

func WorkerNum(workernum int) Option {
	return func(q *Queue) error {
		q.workernum = workernum
		return nil
	}
}

func WithEventbus(eb eventbus.Bus) Option {
	return func(q *Queue) error {
		q.eventbus = eb
		return nil
	}
}

func WithLogger(l logger.Logger) Option {
	return func(q *Queue) error {
		q.logger = l
		return nil
	}
}

// Run starts the queue
func (q *Queue) Run() {
	defer func() {
		q.fire(event.NewQueueAfterStop(q.Name(), q.Uptime()))
	}()
	q.fire(event.NewQueueBeforeRun(q.Name()))
	q.startedAt = time.Now()
	if q.workernum <= 0 {
		q.workernum = DefaultWorkerNum
	}
	q.workers = make(chan queue.Worker, q.workernum)
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
