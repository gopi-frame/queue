package queue

import (
	"time"

	"github.com/google/uuid"
	"github.com/gopi-frame/future"
	"github.com/gopi-frame/queue/driver"
)

func NewWorker(queue driver.Queue) *Worker {
	return &Worker{
		id: uuid.New(),
	}
}

type Worker struct {
	id     uuid.UUID
	queue  driver.Queue
	booted bool
	Quit   chan struct{}
}

func (w *Worker) Start() {
	if w.booted {
		return
	}
	future.Void(func() {
		w.booted = true
		for {
			select {
			case <-w.Quit:
				return
			default:
				job := w.queue.Dequeue()
				if job == nil {
					time.Sleep(time.Second * 5)
					continue
				}
				if err := job.Handle(); err != nil {
					if job.GetModel().GetAttempts() >= job.GetMaxAttempts() {
						job.Failed(err)
					} else {
						w.queue.Release(job, job.GetRetryDelay())
					}
					continue
				}
				w.queue.Ack(job)
			}
		}
	}).Complete(func() { w.booted = false })
}

func (w *Worker) Stop() {
	future.Void(func() { w.Quit <- struct{}{} })
}
