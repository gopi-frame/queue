package queue

import (
	"fmt"
	"sort"

	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/queue/driver"
)

var drivers = kv.NewMap[string, driver.Driver]()

func Register(driverName string, driver driver.Driver) {
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

func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	list := drivers.Keys()
	sort.Strings(list)
	return list
}

func Open(driverName string, options map[string]any) (driver.Queue, error) {
	drivers.RLock()
	driver, ok := drivers.Get(driverName)
	drivers.RUnlock()
	if !ok {
		return nil, exception.NewArgumentException("driverName", driverName, fmt.Sprintf("unknown driver \"%s\"", driverName))
	}
	return driver.Open(options)
}
