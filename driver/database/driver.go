package database

import (
	"github.com/gopi-frame/contract/queue"
	queueimpl "github.com/gopi-frame/queue"
)

var driverName = "database"

func init() {
	if driverName != "" {
		queueimpl.Register(driverName, new(Driver))
	}
}

type Driver struct {
}

func (d Driver) Open(options map[string]any) (queue.Queue, error) {
	cfg, err := ParseCfg(options)
	if err != nil {
		return nil, err
	}
	return NewQueue(cfg), nil
}
