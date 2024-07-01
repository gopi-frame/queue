package driver

import "time"

type Job interface {
	SetModel(Queueable)
	GetModel() Queueable
	GetDelay() time.Duration
	GetMaxAttempts() int
	GetRetryDelay() time.Duration
	Handle() error
	Failed(error)
}
