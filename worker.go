package queue

import (
	"github.com/gopi-frame/contract/eventbus"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue/event"
)

// Worker is a worker implementation
type Worker struct {
	workers  chan queue.Worker
	queue    queue.Queue
	eventbus eventbus.Bus
}

// NewWorker creates a new worker
func NewWorker(workers chan queue.Worker, queue queue.Queue, eventbus eventbus.Bus) *Worker {
	return &Worker{workers: workers, queue: queue, eventbus: eventbus}
}

// Handle handles a job
func (w *Worker) Handle(job queue.Job) {
	defer func() {
		w.workers <- w
	}()
	model := job.GetQueueable()
	w.fire(event.NewJobBeforeHandle(model.GetID(), model.GetQueue()))
	err := job.Handle()
	if err == nil {
		w.fire(event.NewJobAfterHandle(model.GetID(), model.GetQueue()))
		return
	}
	if model.GetAttempts() < job.GetMaxAttempts() {
		w.queue.Release(job)
		w.fire(event.NewJobAfterRelease(model.GetQueue(), err, model.GetID(), model.GetAttempts()+1))
		return
	}
	job.Failed(err)
}

func (w *Worker) fire(e eventbus.Event) {
	if w.eventbus != nil {
		err := w.eventbus.Dispatch(e)
		if err != nil {
			return
		}
	}
}
