package empty

import (
	"github.com/google/uuid"
	queuecontract "github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/contract/support"
	"github.com/gopi-frame/support/lists"
)

// Failer empty failer
type Failer struct{}

// Save save
func (f *Failer) Save(queue string, job queuecontract.Job, err error) {}

// All all
func (f *Failer) All(queue string) support.List[queuecontract.Job] {
	return lists.NewList[queuecontract.Job]()
}

// Find find
func (f *Failer) Find(queue string, id uuid.UUID) queuecontract.Job {
	return nil
}

// Forget forget
func (f *Failer) Forget(queue string, id uuid.UUID) {}

// Flush flush
func (f *Failer) Flush(queue string) {}
