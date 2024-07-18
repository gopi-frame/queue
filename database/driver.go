package database

import (
	qc "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue"
)

var driverName = "database"

func init() {
	if driverName != "" {
		queue.Register(driverName, new(Driver))
	}
}

type Driver struct{}

func (d Driver) Open(options map[string]any) (qc.Queue, error) {
	cfg, err := ParseCfg(options)
	if err != nil {
		return nil, err
	}
	return NewQueue(cfg), nil
}
