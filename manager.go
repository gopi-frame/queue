package queue

import (
	"errors"

	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/support/maps"
)

// NewManager creates a new workerpool manager
func NewManager() *Manager {
	manager := new(Manager)
	manager.queues = maps.NewMap[string, *Queue]()
	return manager
}

// Manager workerpool manager
type Manager struct {
	queues *maps.Map[string, *Queue]
}

// List lists all registered worker pools
func (wpm *Manager) List() (items map[string]*Queue) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	items = wpm.queues.ToMap()
	return
}

// Get returns Worker pool by the specific name
func (wpm *Manager) Get(name string) (dispatcher *Queue) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	dispatcher = wpm.queues.GetOr(name, nil)
	return
}

// Create creates a new worker pool with max worker count and registers it with the specific name
func (wpm *Manager) Create(queue queue.Queue, options ...Option) (d *Queue, isNew bool) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	if wpm.queues.ContainsKey(queue.Name()) {
		d, isNew = wpm.queues.Get(queue.Name())
		return
	}
	d = New(queue, options...)
	wpm.queues.Set(queue.Name(), d)
	isNew = true
	return
}

// Add registers an existing worker pool
func (wpm *Manager) Add(d *Queue) (success bool, err error) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	if wpm.queues.ContainsKey(d.Name()) {
		success = false
		err = errors.New("exists")
		return
	}
	wpm.queues.Set(d.Name(), d)
	success = true
	return
}

// Start starts the specific worker pool
func (wpm *Manager) Start(name string) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	if workerPool, ok := wpm.queues.Get(name); ok {
		if workerPool.Stopped() {
			workerPool.Start()
		}
	}
}

// Stop stops the specific worker pool
func (wpm *Manager) Stop(name string) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	if workerPool, ok := wpm.queues.Get(name); ok {
		if workerPool.Running() {
			workerPool.Stop()
		}
	}
}

// Release releases the specific worker pool
func (wpm *Manager) Release(name string) {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	if workerPool, ok := wpm.queues.Get(name); ok {
		if workerPool.Running() {
			workerPool.Release()
		}
		wpm.queues.Remove(name)
	}
}

// StartAll starts all worker pools
func (wpm *Manager) StartAll() {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	wpm.queues.Each(func(key string, value *Queue) bool {
		if value.Stopped() {
			go value.Start()
		}
		return true
	})
}

// StopAll stops all worker pools
func (wpm *Manager) StopAll() {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	wpm.queues.Each(func(key string, value *Queue) bool {
		if value.Running() {
			go value.Stop()
		}
		return true
	})
}

// ReleaseAll releases all worker pools
func (wpm *Manager) ReleaseAll() {
	wpm.queues.Lock()
	defer wpm.queues.Unlock()
	wpm.queues.Each(func(key string, value *Queue) bool {
		go value.Release()
		return true
	})
	wpm.queues.Clear()
}
