package database

import (
	"github.com/gopi-frame/queue"
	"github.com/gopi-frame/queue/driver"
)

var driverName = "database"

func init() {
	if driverName != "" {
		queue.Register(driverName, new(Driver))
	}
}

type Driver struct{}

func (d Driver) Open(options map[string]any) (driver.Queue, error) {
	cfg, err := ParseCfg(options)
	if err != nil {
		return nil, err
	}
	return NewQueue(cfg), nil
}
