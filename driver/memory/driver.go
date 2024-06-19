package memory

import (
	"github.com/gopi-frame/contract/queue"
	queueimpl "github.com/gopi-frame/queue"
)

var driverName = "memory"

func init() {
	if driverName != "" {
		queueimpl.Register(driverName, new(Driver))
	}
}

type Driver struct {
}

func (d Driver) Open(options map[string]any) (queue.Queue, error) {
	if name, ok := options[OptKeyName]; ok {
		return NewQueue(name.(string)), nil
	} else {
		return nil, ErrMissingOptionName
	}
}
