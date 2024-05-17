package failed

import (
	uuidlib "github.com/google/uuid"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/support/lists"
)

// MemoryFailedJobProvider failed job provider based on memory
type MemoryFailedJobProvider struct {
	queue string
	items *lists.List[*MemoryFailedJob]
}

// SetQueue set queue name
func (provider *MemoryFailedJobProvider) SetQueue(queue string) {
	provider.queue = queue
}

// Save save failed job
func (provider *MemoryFailedJobProvider) Save(item queue.Queueable, err error) {
	provider.items.Push(NewMemoryFailedJob(provider.queue, item, err))
}

// All get all failed job
func (provider *MemoryFailedJobProvider) All() *lists.List[queue.Queueable] {
	models := lists.NewList[queue.Queueable]()
	provider.items.Each(func(_ int, item *MemoryFailedJob) bool {
		models.Push(item.Payload)
		return true
	})
	return models
}

// Find find failed job by uuid
func (provider *MemoryFailedJobProvider) Find(uuid string) queue.Queueable {
	item, ok := provider.items.FirstWhere(func(value *MemoryFailedJob) bool {
		return value.Payload.UUID() == uuidlib.MustParse(uuid)
	})
	if !ok {
		return nil
	}
	return item.Payload
}

// Forget remove failed job by uuid
func (provider *MemoryFailedJobProvider) Forget(uuid string) {
	provider.items.RemoveWhere(func(value *MemoryFailedJob) bool {
		return value.Payload.UUID() == uuidlib.MustParse(uuid)
	})
}

// Flush clear all failed job
func (provider *MemoryFailedJobProvider) Flush() {
	provider.items.Clear()
}
