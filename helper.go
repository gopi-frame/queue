package queue

import (
	"fmt"

	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/support/maps"
)

var resolvers = maps.NewMap[string, queue.DispatcherResolver]()
var connectors = maps.NewMap[string, queue.DispatcherConnector]()

func RegisterResolver(driver string, resolver queue.DispatcherResolver) {
	resolvers.Lock()
	defer resolvers.Unlock()
	resolvers.Set(driver, resolver)
}

func ResolverFor(driver string) queue.DispatcherResolver {
	resolvers.RLock()
	defer resolvers.RUnlock()
	resolver := resolvers.GetOr(driver, nil)
	if resolver == nil {
		panic(exception.NewArgumentException("driver", driver, fmt.Sprintf("No resolver for driver [%s]", driver)))
	}
	return resolver
}

func Resolve(driver string, config map[string]any) queue.Dispatcher {
	return ResolverFor(driver).Resolve(config)
}
