package queue

import (
	"errors"
	queuecontract "github.com/gopi-frame/contract/queue"
)

func init() {
	Register("test", new(mockDriver))
}

// mockDriver is a mock driver for testing
type mockDriver struct {
}

func (mockDriver) Open(options map[string]any) (queuecontract.Queue, error) {
	name := options["name"].(string)
	if name == "exception" {
		return nil, errors.New("exception")
	}
	q := &mockQueue{
		name: name,
	}
	return q, nil
}
