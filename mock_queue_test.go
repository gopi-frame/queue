package queue

import (
	"github.com/gopi-frame/contract/queue"
	"sync"
	"time"
)

// mockQueue is a mock queue implementation for testing
type mockQueue struct {
	name string
	jobs []*mockQueueableJob
	mu   sync.RWMutex
}

func (m *mockQueue) Name() string {
	return m.name
}

func (m *mockQueue) Empty() bool {
	m.mu.RLock()
	defer m.mu.Unlock()
	return len(m.jobs) == 0
}

func (m *mockQueue) Count() int64 {
	m.mu.RLock()
	defer m.mu.Unlock()
	return int64(len(m.jobs))
}

func (m *mockQueue) Enqueue(job queue.Job) (queue.Job, bool) {
	time.Sleep(job.GetDelay())
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs = append(m.jobs, newMockQueueableJob(m.name, job))
	return job, true
}

func (m *mockQueue) Dequeue() (queue.Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.jobs) == 0 {
		return nil, false
	}
	job := m.jobs[0]
	m.jobs = m.jobs[1:]
	return job.payload, true
}

func (m *mockQueue) Remove(job queue.Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, j := range m.jobs {
		if j.GetID() == job.GetQueueable().GetID() {
			m.jobs = append(m.jobs[:i], m.jobs[i+1:]...)
			return
		}
	}
}

func (m *mockQueue) Ack(_ queue.Job) {}

func (m *mockQueue) Release(job queue.Job) {
	time.Sleep(job.GetRetryDelay())
	m.mu.Lock()
	defer m.mu.Unlock()
	job.GetQueueable().(*mockQueueableJob).attempts++
	m.jobs = append(m.jobs, job.GetQueueable().(*mockQueueableJob))
}
