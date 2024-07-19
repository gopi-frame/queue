package event

type JobFailed struct {
}

func (JobFailed) Topic() string {
	return jobFailed
}
