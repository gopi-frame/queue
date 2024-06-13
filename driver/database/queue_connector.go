package database

import "github.com/gopi-frame/contract/queue"

type QueueConnector struct {
}

func (q *QueueConnector) Connect(options ...queue.DispatcherOption) {
	panic("not implemented") // TODO: Implement
}
