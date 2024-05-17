package failed

import (
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/contract/support"
	"github.com/gopi-frame/support/lists"
)

// EmptyFailedJobProvider empty failed job provider
type EmptyFailedJobProvider struct {
}

// Save log
func (provider *EmptyFailedJobProvider) Save(queue string, job queue.Job, err error) {
}

// All all
func (provider *EmptyFailedJobProvider) All() support.List[queue.Job] {
	return lists.NewList[queue.Job]()
}

// Find find
func (provider *EmptyFailedJobProvider) Find(uuid string) queue.Job {
	return nil
}

// Forget forget
func (provider *EmptyFailedJobProvider) Forget(uuid string) {
}

// Flush flush
func (provider *EmptyFailedJobProvider) Flush() {
}
