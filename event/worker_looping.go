package event

import "github.com/google/uuid"

type WorkerLooping struct {
	id    uuid.UUID
	queue string
}

func (WorkerLooping) Topic() string {
	return workerLooping
}

func (w *WorkerLooping) ID() uuid.UUID {
	return w.id
}

func (w *WorkerLooping) Queue() string {
	return w.queue
}

func NewWorkerLooping(id uuid.UUID, queue string) *WorkerLooping {
	return &WorkerLooping{
		id,
		queue,
	}
}
