package job

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

// NewMemoryJob create memory job
func NewMemoryJob(payload queue.Job) *MemoryJob {
	return &MemoryJob{
		uuid:      uuid.New(),
		executing: false,
		payload:   payload,
		attempts:  0,
	}
}

// MemoryJob memory job
type MemoryJob struct {
	uuid      uuid.UUID
	executing bool
	payload   queue.Job
	attempts  uint
}

// ID return id
func (m *MemoryJob) ID() uint64 {
	return 0
}

// UUID return uuid
func (m *MemoryJob) UUID() uuid.UUID {
	return m.uuid
}

// Attempts return attempts
func (m *MemoryJob) Attempts() uint {
	return m.attempts
}

// Payload return payload
func (m *MemoryJob) Payload() queue.Job {
	return m.payload
}

// Executing return executing
func (m *MemoryJob) Executing() bool {
	return m.executing
}
