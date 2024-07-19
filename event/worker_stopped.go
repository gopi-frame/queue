package event

import "github.com/google/uuid"

type WorkerStopped struct {
	id    uuid.UUID
	queue string
}

func (WorkerStopped) Topic() string {
	return workerStopped
}

func (w *WorkerStopped) ID() uuid.UUID {
	return w.id
}

func (w *WorkerStopped) Queue() string {
	return w.queue
}

func NewWorkerStopped(id uuid.UUID, queue string) *WorkerStopped {
	return &WorkerStopped{
		id,
		queue,
	}
}
