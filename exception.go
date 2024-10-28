package queue

import (
	"fmt"
	. "github.com/gopi-frame/contract/exception"
	"github.com/gopi-frame/exception"
)

type QueueNotConfiguredException struct {
	queue string
	Throwable
}

func NewQueueNotConfiguredException(queue string) *QueueNotConfiguredException {
	return &QueueNotConfiguredException{
		queue:     queue,
		Throwable: exception.New(fmt.Sprintf("queue [%s] not configured", queue)),
	}
}

func (err *QueueNotConfiguredException) Queue() string {
	return err.queue
}
