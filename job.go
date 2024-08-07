package queue

import (
	"github.com/gopi-frame/contract/queue"
	"time"
)

// Job is an abstract queue job.
//
// Example:
//
//	# Job With default Configs
//
//	type MyJob struct {
//		queue.Job
//		// Add any additional fields here
//	}
//
//	func (j *MyJob) Handle() error {
//		// Do something
//		return nil
//	}
//
//	func (j *MyJob) Failed(err error) {
//		// Handle job failure
//	}
//
//	# Job With custom Configs
//
//	type MyJob struct {
//		queue.Job
//		// Add any additional fields here
//	}
//
//	func (j *MyJob) Handle() error {
//		// Do something
//		return nil
//	}
//
//	func (j *MyJob) Failed(err error) {
//		// Handle job failure
//	}
//
//	// Set the delay for the job
//	func (j *MyJob) GetDelay() time.Duration {
//		return time.Second * 30
//	}
type Job struct {
	queue.Job
	model queue.Queueable
}

func (j *Job) SetQueueable(model queue.Queueable) {
	j.model = model
}

func (j *Job) GetQueueable() queue.Queueable {
	return j.model
}

func (j *Job) GetDelay() time.Duration {
	return 0
}

func (j *Job) GetMaxAttempts() int {
	return 3
}

func (j *Job) GetRetryDelay() time.Duration {
	return 0
}
