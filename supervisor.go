package queue

import (
	"fmt"

	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/support/maps"
)

var queueConnectors = maps.NewMap[string, queue.DispatcherConnector]()

func RegisterQueueConnector(driver string, connector queue.DispatcherConnector) {
	queueConnectors.Lock()
	defer queueConnectors.Unlock()
	queueConnectors.Set(driver, connector)
}

func QueueConnectors() []string {
	queueConnectors.RLock()
	defer queueConnectors.RUnlock()
	return queueConnectors.Keys()
}

func QueueConnectorFor(driver string) queue.DispatcherConnector {
	queueConnectors.RLock()
	defer queueConnectors.RUnlock()
	connector, ok := queueConnectors.Get(driver)
	if !ok {
		panic(exception.NewArgumentException("driver", driver, fmt.Sprintf("No connector for driver [%s]", driver)))
	}
	return connector
}

func OpenQueue(driver string, config map[string]any) queue.Dispatcher {
	return QueueConnectorFor(driver).Connect()
}

var failedJobProviderConnectors = maps.NewMap[string, queue.FailedJobProviderResolver]()

func RegisterFailedJobProviderConnector(driver string, connector queue.FailedJobProviderResolver) {
	failedJobProviderConnectors.Lock()
	defer failedJobProviderConnectors.Unlock()
	failedJobProviderConnectors.Set(driver, connector)
}

func FailedJobProviderConnectors() []string {
	failedJobProviderConnectors.RLock()
	defer failedJobProviderConnectors.RUnlock()
	return failedJobProviderConnectors.Keys()
}

func FailedJobProviderConnectorFor(driver string) queue.FailedJobProviderResolver {
	connector, ok := failedJobProviderConnectors.Get(driver)
	if ok {
		panic(exception.NewArgumentException("driver", driver, fmt.Sprintf("No connector for driver [%s]", driver)))
	}
	return connector
}

func OpenFailedJobProvider(driver string, config map[string]any) queue.FailedJobProvider {
	return FailedJobProviderConnectorFor(driver).Connect(config)
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		queues:    maps.NewMap[string, queue.Dispatcher](),
		concrates: maps.NewMap[string, func() queue.Dispatcher](),
	}
}

type Supervisor struct {
	defaultQueue string
	queues       *maps.Map[string, queue.Dispatcher]
	concrates    *maps.Map[string, func() queue.Dispatcher]
	failed       queue.FailedJobProvider
}

func (s *Supervisor) SetDefaultQueue(queue string) {
	s.defaultQueue = queue
}

func (s *Supervisor) SetFailedJobProvider(provider queue.FailedJobProvider) {
	s.failed = provider
}

func (s *Supervisor) Add(name string, queue queue.Dispatcher) {
	s.queues.Lock()
	defer s.queues.Unlock()
	s.queues.Set(name, queue)
}

func (s *Supervisor) Bind(name string, concrate func() queue.Dispatcher) {
	s.concrates.Lock()
	defer s.concrates.Unlock()
	s.concrates.Set(name, concrate)
}

func (s *Supervisor) Dispatch(task queue.Job) bool {
	return s.DispatchTo(task, s.defaultQueue)
}

func (s *Supervisor) DispatchTo(task queue.Job, name string) bool {
	s.queues.RLock()
	defer s.queues.RUnlock()
	queue := s.queues.GetOr(name, nil)
	if queue == nil {
		s.concrates.RLock()
		defer s.concrates.RUnlock()
		if concrate := s.concrates.GetOr(name, nil); concrate == nil {
			return false
		} else {
			queue := concrate()
			s.queues.Set(name, queue)
			return queue.Dispatch(task)
		}
	}
	return queue.Dispatch(task)
}
