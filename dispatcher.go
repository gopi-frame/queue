package queue

import (
	"time"

	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/future"
)

const DefaultNumProcs = 3

func NewDispatcher(queue queue.Queue, numprocs uint) queue.Dispatcher {
	q := &Dispatcher{
		queue:    queue,
		numprocs: numprocs,
	}
	return q
}

type Dispatcher struct {
	name     string
	queue    queue.Queue
	numprocs uint
	failed   queue.FailedJobProvider
}

func (w *Dispatcher) FailedJobProvider(provider queue.FailedJobProvider) {
	w.failed = provider
}

func (w *Dispatcher) Dispatch(job queue.Job) error {
	return w.queue.Enqueue(job)
}

func (w *Dispatcher) Reload() {
	if w.failed == nil {
		return
	}
	failedJobs := w.failed.All(w.name)
	for _, failedJob := range failedJobs {
		w.Dispatch(failedJob.GetPayload())
	}
}

func (w *Dispatcher) Flush() {
	if w.failed == nil {
		return
	}
	w.failed.Flush(w.name)
}

func (w *Dispatcher) Exec() {
	for i := 0; i < int(w.numprocs); i++ {
		go func() {
			for {
				job, err := w.queue.Dequeue()
				if !ok {
					time.Sleep(time.Second * 5)
					continue
				}
				future.Void(func() {
					if err := job.Handle(); err != nil {
						panic(err)
					}
				}).CatchAll(job.Failed).Await()
				time.Sleep(time.Second)
			}
		}()
	}
}
