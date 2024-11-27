package queue

import (
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/queue"
)

type QueueManager struct {
	queues *kv.Map[string, queue.Queue]
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: kv.NewMap[string, queue.Queue](),
	}
}

func (qm *QueueManager) AddQueue(name string, queue queue.Queue) {
	qm.queues.Lock()
	defer qm.queues.Unlock()
	qm.queues.Set(name, queue)
}

func (qm *QueueManager) HasQueue(name string) bool {
	qm.queues.RLock()
	if qm.queues.ContainsKey(name) {
		qm.queues.RUnlock()
		return true
	}
	qm.queues.RUnlock()
	return false
}

func (qm *QueueManager) TryGetQueue(name string) (queue.Queue, error) {
	qm.queues.RLock()
	if q, ok := qm.queues.Get(name); ok {
		qm.queues.RUnlock()
		return q, nil
	}
	qm.queues.RUnlock()
	return nil, NewQueueNotConfiguredException(name)
}

func (qm *QueueManager) GetQueue(name string) queue.Queue {
	if q, err := qm.TryGetQueue(name); err != nil {
		panic(err)
	} else {
		return q
	}
}
