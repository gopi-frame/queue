package queue

import (
	"sort"

	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/support/maps"
)

var drivers = maps.NewMap[string, queue.Driver]()

func Register(name string, driver queue.Driver) {
	drivers.Lock()
	defer drivers.Unlock()
	if driver == nil {
		panic(exception.NewEmptyArgumentException("driver"))
	}
	if _, dup := drivers.Get(name); dup {
		panic(NewDuplicateDriverException(name))
	}
	drivers.Set(name, driver)
}

func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	list := drivers.Keys()
	sort.Strings(list)
	return list
}

func Open(driverName string, options map[string]any) (queue.Queue, error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		panic(NewUnknownDriverException(driverName))
	}
	return driver.Open(options)
}
