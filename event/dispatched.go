package event

type Dispatched struct {
	queue string
}

func (Dispatched) Topic() string {
	return dispatched
}

func (d *Dispatched) Queue() string {
	return d.queue
}

func NewDispatched(queue string) *Dispatched {
	return &Dispatched{
		queue,
	}
}
