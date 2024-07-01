package memory

import (
	queueimpl "github.com/gopi-frame/queue"
	"github.com/gopi-frame/queue/driver"
)

var driverName = "memory"

func init() {
	if driverName != "" {
		queueimpl.Register(driverName, new(Driver))
	}
}

type Driver struct {
}

func (d Driver) Open(options map[string]any) (driver.Queue, error) {
	if name, ok := options[OptKeyName]; ok {
		return NewQueue(name.(string)), nil
	} else {
		return nil, ErrMissingOptionName
	}
}
