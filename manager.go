package queue

import (
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/future"
	"github.com/gopi-frame/support/maps"
)

func NewManager() *Manager {
	return &Manager{
		dispatchers: maps.NewMap[string, *Dispatcher](),
	}
}

type Manager struct {
	dispatchers *maps.Map[string, *Dispatcher]
}

func (m *Manager) Set(queue string, dispatcher *Dispatcher) {
	if m.dispatchers.TryLock() {
		defer m.dispatchers.Unlock()
	}
	m.dispatchers.Set(queue, dispatcher)
}

func (m *Manager) Get(queue string) *Dispatcher {
	if m.dispatchers.TryLock() {
		defer m.dispatchers.Unlock()
	}
	return m.dispatchers.GetOr(queue, nil)
}

func (m *Manager) DispatchTo(queue string, job queue.JobInterface) {
	if dispatcher := m.Get(queue); dispatcher != nil {
		dispatcher.Dispatch(job)
	}
}

func (m *Manager) Start() {
	m.dispatchers.Each(func(queue string, dispatcher *Dispatcher) bool {
		future.Void(dispatcher.Start)
		return true
	})
}
