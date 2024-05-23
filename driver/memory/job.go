package memory

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
)

// Job job implemention
type Job struct {
	uuid     uuid.UUID
	payload  queue.Job
	attempts uint
}

// ID id
func (j *Job) ID() uint64 {
	return 0
}

// UUID uuid
func (j *Job) UUID() uuid.UUID {
	return j.uuid
}

// Attempts attempts
func (j *Job) Attempts() uint {
	return j.attempts
}

// Payload payload
func (j *Job) Payload() queue.Job {
	return j.payload
}
