// Package redis is a redis queue driver for queue package.
package redis

import (
	"github.com/go-viper/mapstructure/v2"
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/queue"
)

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/queue/driver/redis.driverName=custom`
var driverName = "redis"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		queue.Register(driverName, new(Driver))
	}
}

// Driver is the redis queue driver
type Driver struct {
}

// Open opens the redis queue
func (Driver) Open(options map[string]any) (queuecontract.Queue, error) {
	cfg := new(Config)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           cfg,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(options); err != nil {
		return nil, err
	}
	return NewQueue(cfg), nil
}

// Open is a convenience function that calls [Driver.Open].
func Open(options map[string]any) (queuecontract.Queue, error) {
	return new(Driver).Open(options)
}
