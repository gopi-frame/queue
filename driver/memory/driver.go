package memory

import (
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/queue"
)

// This variable can be replaced through `go build -ldflags=-X github.com/gopi-frame/queue/driver/memory.driverName=custom`
var driverName = "memory"

//goland:noinspection GoBoolExpressions
func init() {
	if driverName != "" {
		queue.Register(driverName, new(Driver))
	}
}

// Driver provides memory queue driver
type Driver struct {
}

// Open opens memory queue, options must contain "name"
func (d Driver) Open(options map[string]any) (queuecontract.Queue, error) {
	if name, ok := options["name"]; ok && name.(string) != "" {
		return NewQueue(name.(string)), nil
	} else {
		return nil, exception.New("options[\"name\"] is required")
	}
}

// Open is a convenience function that calls [Driver.Open].
func Open(options map[string]any) (queuecontract.Queue, error) {
	return new(Driver).Open(options)
}
