package event

import (
	"github.com/gopi-frame/contract/event"

	events "github.com/gopi-frame/event"
)

// Subscriber subscriber
type Subscriber struct {
	BeforeHandle    func(event.Event) bool
	AfterHandle     func(event.Event) bool
	FailedHandle    func(event.Event) bool
	ProgressUpdated func(event.Event) bool
}

// OnBeforeHandle handles on before event
func (subscriber *Subscriber) OnBeforeHandle(event event.Event) bool {
	if subscriber.BeforeHandle != nil {
		return subscriber.BeforeHandle(event)
	}
	return true
}

// OnAfterHandle handles on after event
func (subscriber *Subscriber) OnAfterHandle(event event.Event) bool {
	if subscriber.AfterHandle != nil {
		return subscriber.AfterHandle(event)
	}
	return true
}

// OnFailedHandle handles on failed event
func (subscriber *Subscriber) OnFailedHandle(event event.Event) bool {
	if subscriber.FailedHandle != nil {
		return subscriber.FailedHandle(event)
	}
	return true
}

// OnProgressUpdated on progress updated
func (subscriber *Subscriber) OnProgressUpdated(event event.Event) bool {
	if subscriber.ProgressUpdated != nil {
		return subscriber.ProgressUpdated(event)
	}
	return true
}

// Subscribe returns top-event map
func (subscriber *Subscriber) Subscribe(publisher event.Dispatcher) {
	publisher.Listen([]event.Event{new(BeforeHandle)}, events.Listener(subscriber.OnBeforeHandle))
	publisher.Listen([]event.Event{new(AfterHandle)}, events.Listener(subscriber.OnAfterHandle))
	publisher.Listen([]event.Event{new(Failed)}, events.Listener(subscriber.OnFailedHandle))
	publisher.Listen([]event.Event{new(ProgressUpdated)}, events.Listener(subscriber.OnProgressUpdated))
}
