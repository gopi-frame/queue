package driver

type Queueable interface {
	GetID() string
	GetQueue() string
	GetPayload() Job
	GetAttempts() int
}
