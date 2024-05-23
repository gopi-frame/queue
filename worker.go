package queue

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/future"
	"github.com/gopi-frame/queue/event"
)

// WorkerStatus worker status
type WorkerStatus int

// worker status enums
const (
	WorkerStatusIdle WorkerStatus = iota + 1
	WorkerStatusWorking
	WorkerStatusStopped
)

func newWorker(queue *Queue) *Worker {
	worker := Worker{
		id:          uuid.New(),
		status:      WorkerStatusIdle,
		createdAt:   time.Now(),
		idledAt:     time.Now(),
		queue:       queue,
		stopChannel: make(chan struct{}),
	}
	return &worker
}

// Worker is a struct to handle jobs
type Worker struct {
	id          uuid.UUID    // unique id
	status      WorkerStatus // status
	createdAt   time.Time    // created time
	startedAt   time.Time    // last started time
	idledAt     time.Time    // last idled time
	stoppedAt   time.Time    // last stopped time
	queue       *Queue
	stopChannel chan struct{}
}

func (worker *Worker) setWorking() {
	worker.status = WorkerStatusWorking
	worker.startedAt = time.Now()
}

func (worker *Worker) setIdle() {
	worker.status = WorkerStatusIdle
	worker.idledAt = time.Now()
}

func (worker *Worker) setStopped() {
	worker.status = WorkerStatusStopped
	worker.stoppedAt = time.Now()
}

// ID returns worker's unique id
func (worker *Worker) ID() uuid.UUID {
	return worker.id
}

// Status returns worker's active status
//   - [WorkerStatusIdle]
//   - [WorkerStatusWorking]
//   - [WorkerStatusStopped]
func (worker *Worker) Status() WorkerStatus {
	return worker.status
}

// CreatedAt returns the worker's created time
func (worker *Worker) CreatedAt() time.Time {
	return worker.createdAt
}

// IdledAt returns the worker's last idle time
func (worker *Worker) IdledAt() time.Time {
	return worker.idledAt
}

// StoppedAt returns the worker's last stopped time
func (worker *Worker) StoppedAt() time.Time {
	return worker.stoppedAt
}

// Working returns whether the worker is working
func (worker *Worker) Working() bool {
	return worker.status == WorkerStatusWorking
}

// Idle returns whether the worker is idle
func (worker *Worker) Idle() bool {
	return worker.status == WorkerStatusIdle
}

// Stopped returns whether the worker is stopped
func (worker *Worker) Stopped() bool {
	return worker.status == WorkerStatusStopped
}

// Stoppable returns whether the worker can be stopped
func (worker *Worker) Stoppable() bool {
	return worker.status == WorkerStatusIdle || worker.status == WorkerStatusWorking
}

func (worker *Worker) handle(job queue.Job) {
	worker.setWorking()
	timeout := job.Timeout()
	if worker.queue != nil {
		worker.queue.DispatchEvent(event.NewBeforeHandle(job))
	}
	future.Timeout(func() error {
		return job.Handle()
	}, timeout).Then(func(value error) error {
		if value == nil {
			if worker.queue != nil {
				worker.queue.Ack(job)
				worker.queue.DispatchEvent(event.NewAfterHandle(job))
				pending, executing := worker.queue.Progress()
				worker.queue.DispatchEvent(event.NewProgressUpdated(pending, executing))
			}
		}
		return nil
	}, nil).CatchAll(func(err error) {
		if worker.queue != nil {
			worker.queue.Fail(job, err)
			if job.MaxAttempts() <= job.GetJob().Attempts() {
				worker.queue.DispatchEvent(event.NewFailed(job, err))
				pending, executing := worker.queue.Progress()
				worker.queue.DispatchEvent(event.NewProgressUpdated(pending, executing))
			}
		}
	}).Await()
	worker.setIdle()
}

// Start lets the worker starting worker
func (worker *Worker) Start() {
	go func() {
		for {
			if worker.ShouldStop() {
				worker.Stop()
				continue
			}
			if worker.ShouldRelease() {
				worker.Release()
				return
			}
			time.Sleep(time.Second)
		}
	}()
	for {
		select {
		case <-worker.stopChannel:
			worker.setStopped()
			return
		default:
			if job, ok := worker.queue.Dequeue(); ok {
				worker.handle(job)
			} else {
				time.Sleep(time.Second)
			}
		}
	}
}

// Stop stops the worker, if the worker's status is [WorkerStatusIdle] it will be stopped immediately
// if the worker's status is [WorkerStatusWorking], its status will change to [WorkerStatusStopping] and
// will be stopped after MaxExecuteTimePerAttempt
func (worker *Worker) Stop() {
	worker.stopChannel <- struct{}{}
}

// Release releases the worker
func (worker *Worker) Release() {
	defer func() {
		worker.queue.workers.Remove(worker.id)
		worker.queue = nil
	}()
	worker.Stop()
	close(worker.stopChannel)
}

// ShouldStop returns if the worker should be stopped
// it will return true when worker's status is [WorkerStatusIdle] and has been idled over max idle time
func (worker *Worker) ShouldStop() bool {
	return worker.Idle() && time.Since(worker.idledAt) >= worker.queue.workerMaxIdleTime
}

// ShouldRelease returns if the worker should be released,
// it will return true when worker's status is [WorkerStatusStopped] and has been stopped over max stopped time
func (worker *Worker) ShouldRelease() bool {
	return worker.Stopped() && time.Since(worker.stoppedAt) >= worker.queue.workerMaxStoppedTime
}
