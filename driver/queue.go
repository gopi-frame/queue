package driver

import (
	"github.com/gopi-frame/contract/event"
	"github.com/gopi-frame/contract/queue"
)

var _ queue.Queue = (*Queue)(nil)

// Queue abstract driver
type Queue struct {
	queue.Queue
	events event.Dispatcher
}

// Subscribe add a subscriber to queue events
func (driver *Queue) Subscribe(subscriber queue.Subcriber) {
	if driver.events != nil {
		driver.events.Subscribe(subscriber)
	}
}

// DispatchEvent dispatches specifia event
func (driver *Queue) DispatchEvent(event event.Event) {
	driver.events.Dispatch(event)
}
