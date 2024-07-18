package memory

import (
	qc "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue"
)

var driverName = "memory"

func init() {
	if driverName != "" {
		queue.Register(driverName, new(Driver))
	}
}

type Driver struct {
}

func (d Driver) Open(options map[string]any) (qc.Queue, error) {
	if name, ok := options[OptKeyName]; ok {
		return NewQueue(name.(string)), nil
	} else {
		return nil, ErrMissingOptionName
	}
}
