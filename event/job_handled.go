package event

type JobHandled struct {
}

func (JobHandled) Topic() string {
	return jobHandled
}
