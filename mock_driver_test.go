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
	if name == "error" {
		return nil, errors.New("error")
	}
	q := &mockQueue{
		name: name,
	}
	return q, nil
}
