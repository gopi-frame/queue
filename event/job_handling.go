package event

type JobHandling struct {
}

func (JobHandling) Topic() string {
	return jobHandling
}
