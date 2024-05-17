package dispatcher

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/support/maps"
)

// Status workerpool status
type Status int

// WorkerPoolStatus enums
const (
	Stopped Status = iota
	Running
)

// Default default dispatcher
func Default(queue queue.Queue) *Dispatcher {
	return &Dispatcher{
		Queue:                queue,
		workers:              maps.NewMap[uuid.UUID, *Worker](),
		workerNum:            DefaultMaxWorkerNum,
		workerBatch:          DefaultWorkerCreationBatch,
		workerMaxIdleTime:    DefaultWorkerMaxIdleTime,
		workerMaxStoppedTime: DefaultWorkerMaxStoppedTime,
	}
}

// New new dispatcher
func New(queue queue.Queue, options ...Option) *Dispatcher {
	dispatcher := Default(queue)
	for _, option := range options {
		option(dispatcher)
	}
	return dispatcher
}

// Dispatcher is a struct to manage workers
// it accepts Job and push the job to workers
type Dispatcher struct {
	queue.Queue
	status               Status
	workers              *maps.Map[uuid.UUID, *Worker]
	workerNum            int
	workerBatch          int
	workerMaxIdleTime    time.Duration
	workerMaxStoppedTime time.Duration
	shouldRetry          func(error) bool
}

// Name returns the name of the WorkerPool if it's added into WorkerPoolManager
//
// It will return empty string if this WorkerPool instance is not added into WorkerPoolManager
func (wp *Dispatcher) Name() string {
	return wp.Queue.Name()
}

// Status returns the active status of the WorkerPool
func (wp *Dispatcher) Status() Status {
	return wp.status
}

// Running returns whether the WorkerPool is running
func (wp *Dispatcher) Running() bool {
	return wp.status == Running
}

// Stopped returns whether the WorkerPool is stopped
func (wp *Dispatcher) Stopped() bool {
	return wp.status == Stopped
}

func (wp *Dispatcher) setStopped() {
	wp.status = Stopped
}

func (wp *Dispatcher) setStarted() {
	wp.status = Running
}

// Dispatch dispatches job
func (wp *Dispatcher) Dispatch(job queue.Job) bool {
	if wp.Stopped() {
		return false
	}
	ok := wp.Enqueue(job)
	wp.createWorkers()
	return ok
}

// Start starts the workerpool
func (wp *Dispatcher) Start() {
	wp.setStarted()
	wp.createWorkers()
}

// Stop stops the worker pool and all the workers
func (wp *Dispatcher) Stop() {
	wp.workers.Lock()
	defer wp.workers.Unlock()
	wp.workers.Each(func(_ uuid.UUID, worker *Worker) bool {
		worker.Stop()
		return true
	})
	wp.setStopped()
}

// Release releases and removes the workerpool from the [Manager]
func (wp *Dispatcher) Release() {
	// if the worker pool is running, stop it first
	if wp.Running() {
		wp.setStopped()
	}
	// release workers
	workers := wp.Workers()
	for _, worker := range workers {
		worker.Release()
	}
	wp.workers.Clear()
}

func (wp *Dispatcher) createWorkers() {
	wp.workers.Lock()
	defer wp.workers.Unlock()
	if wp.workers.Count() >= int64(wp.workerNum) {
		return
	}
	if wp.IsEmpty() {
		return
	}
	need := int64(wp.workerNum / wp.workerBatch)
	if need == 0 {
		need = int64(wp.workerNum)
	}
	if count := wp.Count(); need > count {
		need = count
	}
	var c int64 = 0
	// awake sleeping workers
	wp.workers.Each(func(_ uuid.UUID, worker *Worker) bool {
		if worker.Stopped() {
			go worker.Start()
			c++
		}
		return true
	})
	// create new workers
	for i := c; i < need; i++ {
		if wp.workers.Count() >= int64(wp.workerNum) {
			return
		}
		w := newWorker(wp)
		wp.workers.Set(w.id, w)
		go w.Start()
	}
}

// Workers returns a slice of Workers
func (wp *Dispatcher) Workers() (workers []*Worker) {
	wp.workers.Lock()
	defer wp.workers.Unlock()
	workers = wp.workers.Values()
	return
}
