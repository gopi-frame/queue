package queue

import "github.com/gopi-frame/contract/queue"

type DeferQueue struct {
	queue.Queue

	driver string
	config map[string]any
}

func NewDeferQueue(driver string, config map[string]any) *DeferQueue {
	return &DeferQueue{
		driver: driver,
		config: config,
	}
}

func (q *DeferQueue) deferInit() {
	if q.Queue != nil {
		return
	}
	var err error
	if q.Queue, err = Open(q.driver, q.config); err != nil {
		panic(err)
	}
	return
}

func (q *DeferQueue) Name() string {
	q.deferInit()
	return q.Queue.Name()
}

func (q *DeferQueue) Empty() bool {
	q.deferInit()
	return q.Queue.Empty()
}

func (q *DeferQueue) Count() int64 {
	q.deferInit()
	return q.Queue.Count()
}

func (q *DeferQueue) Enqueue(job queue.Job) (queue.Job, bool) {
	q.deferInit()
	return q.Queue.Enqueue(job)
}

func (q *DeferQueue) Dequeue() (queue.Job, bool) {
	q.deferInit()
	return q.Queue.Dequeue()
}

func (q *DeferQueue) Remove(job queue.Job) {
	q.deferInit()
	q.Queue.Remove(job)
}

func (q *DeferQueue) Ack(job queue.Job) {
	q.deferInit()
	q.Queue.Ack(job)
}

func (q *DeferQueue) Release(job queue.Job) {
	q.deferInit()
	q.Queue.Release(job)
}
