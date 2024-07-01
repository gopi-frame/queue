package queue

import (
	"time"

	"github.com/gopi-frame/queue/driver"
)

type Job struct {
	driver.Job
	model driver.Queueable
}

func (j *Job) SetModel(model driver.Queueable) {
	j.model = model
}

func (j *Job) GetModel() driver.Queueable {
	return j.model
}

func (j *Job) Release(delay time.Duration) {
}
