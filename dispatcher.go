package queue

import (
	"github.com/gopi-frame/contract/queue"
)

const DefaultWorkers = 3

func NewDispatcher(queue queue.Queue, nworkers int) *Dispatcher {
	if nworkers <= 0 {
		nworkers = DefaultWorkers
	}
	q := &Dispatcher{
		queue:    queue,
		nworkers: nworkers,
	}
	return q
}

type Dispatcher struct {
	queue    queue.Queue
	nworkers int
	workers  []*Worker
	booted   bool
}

func (d *Dispatcher) Dispatch(job queue.JobInterface) {
	d.queue.Enqueue(job)
}

func (d *Dispatcher) Start() {
	for i := 0; i < d.nworkers; i++ {
		worker := NewWorker(d.queue)
		d.workers = append(d.workers, worker)
		worker.Start()
	}
}

func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
	d.booted = false
}
