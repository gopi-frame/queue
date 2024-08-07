package queue

import (
	"fmt"
	"sort"

	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
)

var drivers = kv.NewMap[string, queue.Driver]()

// Register registers a queue driver
func Register(driverName string, driver queue.Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if driver == nil {
		panic(exception.NewEmptyArgumentException("driver"))
	}
	if drivers.ContainsKey(driverName) {
		panic(exception.NewArgumentException("driverName", driverName, fmt.Sprintf("duplicate driver \"%s\"", driverName)))
	}
	drivers.Set(driverName, driver)
}

// Drivers returns list of registered drivers
func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	list := drivers.Keys()
	sort.Strings(list)
	return list
}

// Open opens a queue
func Open(driverName string, options map[string]any) (queue.Queue, error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown driver \"%s\"", driverName))
	}
	return driver.Open(options)
}
