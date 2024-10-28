package queue

import (
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/queue"
)

type QueueManager struct {
	queues      *kv.Map[string, queue.Queue]
	deferQueues *kv.Map[string, func() (queue.Queue, error)]
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

func (qm *QueueManager) AddDeferQueue(name string, config map[string]any) {
	qm.deferQueues.Lock()
	defer qm.deferQueues.Unlock()
	qm.deferQueues.Set(name, func() (queue.Queue, error) {
		driver := config["driver"].(string)
		return Open(driver, config)
	})
}

func (qm *QueueManager) HasQueue(name string) bool {
	qm.queues.RLock()
	if qm.queues.ContainsKey(name) {
		qm.queues.RUnlock()
		return true
	}
	qm.queues.RUnlock()
	qm.deferQueues.RLock()
	if qm.deferQueues.ContainsKey(name) {
		qm.deferQueues.RUnlock()
		return true
	}
	qm.deferQueues.RUnlock()
	return false
}

func (qm *QueueManager) TryQueue(name string) (queue.Queue, error) {
	qm.queues.RLock()
	if q, ok := qm.queues.Get(name); ok {
		qm.queues.RUnlock()
		return q, nil
	}
	qm.queues.RUnlock()
	qm.deferQueues.RLock()
	if q, ok := qm.deferQueues.Get(name); ok {
		qm.deferQueues.RUnlock()
		if q, err := q(); err != nil {
			return nil, err
		} else {
			qm.queues.Lock()
			defer qm.queues.Unlock()
			qm.queues.Set(name, q)
			return q, nil
		}
	}
	qm.deferQueues.RUnlock()
	return nil, NewQueueNotConfiguredException(name)
}

func (qm *QueueManager) Queue(name string) queue.Queue {
	if q, err := qm.TryQueue(name); err != nil {
		panic(err)
	} else {
		return q
	}
}
