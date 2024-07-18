package queue

import (
	"time"

	"github.com/gopi-frame/contract/queue"
)

type Job struct {
	queue.Job
	model queue.Queueable
}

func (j *Job) SetModel(model queue.Queueable) {
	j.model = model
}

func (j *Job) GetModel() queue.Queueable {
	return j.model
}

func (j *Job) Release(delay time.Duration) {
}
