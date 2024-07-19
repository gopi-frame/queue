package event

import "github.com/google/uuid"

type WorkerStopping struct {
	id    uuid.UUID
	queue string
}

func (WorkerStopping) Topic() string {
	return workerStopping
}

func (w *WorkerStopping) ID() uuid.UUID {
	return w.id
}

func (w *WorkerStopping) Queue() string {
	return w.queue
}

func NewWorkerStopping(id uuid.UUID, queue string) *WorkerStopping {
	return &WorkerStopping{
		id,
		queue,
	}
}
