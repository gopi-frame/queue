package event

type Dispatching struct {
	queue string
}

func (Dispatching) Topic() string {
	return dispatching
}

func (d *Dispatching) Queue() string {
	return d.queue
}

func NewDispatching(queue string) *Dispatching {
	return &Dispatching{
		queue,
	}
}
