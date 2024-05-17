package queue

import (
	"errors"

	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue/dispatcher"
	"github.com/gopi-frame/support/maps"
)

// NewManager creates a new workerpool manager
func NewManager() *Manager {
	manager := new(Manager)
	manager.pools = maps.NewMap[string, *dispatcher.Dispatcher]()
	return manager
}

// Manager workerpool manager
type Manager struct {
	pools *maps.Map[string, *dispatcher.Dispatcher]
}

// List lists all registered worker pools
func (wpm *Manager) List() (items map[string]*dispatcher.Dispatcher) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	items = wpm.pools.ToMap()
	return
}

// Get returns Worker pool by the specific name
func (wpm *Manager) Get(name string) (dispatcher *dispatcher.Dispatcher) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	dispatcher = wpm.pools.GetOr(name, nil)
	return
}

// Create creates a new worker pool with max worker count and registers it with the specific name
func (wpm *Manager) Create(queue queue.Queue, options ...dispatcher.Option) (d *dispatcher.Dispatcher, isNew bool) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	if wpm.pools.ContainsKey(queue.Name()) {
		d, isNew = wpm.pools.Get(queue.Name())
		return
	}
	d = dispatcher.New(queue, options...)
	wpm.pools.Set(queue.Name(), d)
	isNew = true
	return
}

// Add registers an existing worker pool
func (wpm *Manager) Add(d *dispatcher.Dispatcher) (success bool, err error) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	if wpm.pools.ContainsKey(d.Name()) {
		success = false
		err = errors.New("exists")
		return
	}
	wpm.pools.Set(d.Name(), d)
	success = true
	return
}

// Start starts the specific worker pool
func (wpm *Manager) Start(name string) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	if workerPool, ok := wpm.pools.Get(name); ok {
		if workerPool.Stopped() {
			workerPool.Start()
		}
	}
}

// Stop stops the specific worker pool
func (wpm *Manager) Stop(name string) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	if workerPool, ok := wpm.pools.Get(name); ok {
		if workerPool.Running() {
			workerPool.Stop()
		}
	}
}

// Release releases the specific worker pool
func (wpm *Manager) Release(name string) {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	if workerPool, ok := wpm.pools.Get(name); ok {
		if workerPool.Running() {
			workerPool.Release()
		}
		wpm.pools.Remove(name)
	}
}

// StartAll starts all worker pools
func (wpm *Manager) StartAll() {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	wpm.pools.Each(func(key string, value *dispatcher.Dispatcher) bool {
		if value.Stopped() {
			go value.Start()
		}
		return true
	})
}

// StopAll stops all worker pools
func (wpm *Manager) StopAll() {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	wpm.pools.Each(func(key string, value *dispatcher.Dispatcher) bool {
		if value.Running() {
			go value.Stop()
		}
		return true
	})
}

// ReleaseAll releases all worker pools
func (wpm *Manager) ReleaseAll() {
	wpm.pools.Lock()
	defer wpm.pools.Unlock()
	wpm.pools.Each(func(key string, value *dispatcher.Dispatcher) bool {
		go value.Release()
		return true
	})
	wpm.pools.Clear()
}
