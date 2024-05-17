package options

import (
	"github.com/gopi-frame/contract/queue"
)

// Subcriber set events subscriber
func Subcriber(subscriber queue.Subcriber) queue.Option {
	return func(queue queue.Queue) {
		queue.Subscribe(subscriber)
	}
}
