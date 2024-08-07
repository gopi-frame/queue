package queue

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"time"
)

type mockJob struct {
	Job             `json:"-"`
	delay           time.Duration
	retryDelay      time.Duration
	handled         bool
	shouldFailTimes int
	failed          bool
	failedReason    error
}

func (m *mockJob) GetDelay() time.Duration {
	return m.delay
}

func (m *mockJob) GetRetryDelay() time.Duration {
	return m.retryDelay
}

func (m *mockJob) Handle() error {
	if m.shouldFailTimes > 0 {
		m.shouldFailTimes--
		return errors.New("fail")
	}
	m.handled = true
	return nil
}

func (m *mockJob) Failed(err error) {
	m.failed = true
	m.failedReason = err
}

type mockQueueableJob struct {
	id       uuid.UUID
	queue    string
	payload  queue.Job
	attempts int
}

func newMockQueueableJob(queue string, payload queue.Job) *mockQueueableJob {
	j := &mockQueueableJob{
		id:       uuid.New(),
		queue:    queue,
		payload:  payload,
		attempts: 0,
	}
	payload.SetQueueable(j)
	return j
}

func (m *mockQueueableJob) GetID() string {
	return m.id.String()
}

func (m *mockQueueableJob) GetQueue() string {
	return m.queue
}

func (m *mockQueueableJob) GetPayload() queue.Job {
	return m.payload
}

func (m *mockQueueableJob) GetAttempts() int {
	return m.attempts
}
