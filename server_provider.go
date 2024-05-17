package queue

import (
	"reflect"

	"github.com/gopi-frame/contract/container"
	"github.com/gopi-frame/contract/queue"
)

// ServerProvider server provider
type ServerProvider struct{}

// Register register
func (s *ServerProvider) Register(c container.Container) {
	c.Bind(reflect.TypeFor[Manager]().String(), func(c container.Container) any {
		return NewManager()
	})
	c.Alias(reflect.TypeFor[Manager]().String(), "queue")
	c.Alias(reflect.TypeFor[Manager]().String(), reflect.TypeFor[queue.Queue]().String())
}
