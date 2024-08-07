package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Run("duplicate name", func(t *testing.T) {
		assert.Panics(t, func() {
			Register("test", new(mockDriver))
		})
	})

	t.Run("nil driver", func(t *testing.T) {
		assert.Panics(t, func() {
			Register("mock", nil)
		})
	})

	t.Run("success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Register("mock", new(mockDriver))
		})
	})
}

func TestOpen(t *testing.T) {
	t.Run("unregistered driver", func(t *testing.T) {
		q, err := Open("unregistered", nil)
		assert.Error(t, err)
		assert.Nil(t, q)
	})

	t.Run("registered driver", func(t *testing.T) {
		q, err := Open("test", map[string]any{
			"name": "test",
		})
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.IsType(t, new(mockQueue), q)
		}
	})

	t.Run("has error", func(t *testing.T) {
		q, err := Open("test", map[string]any{
			"name": "error",
		})
		assert.Error(t, err)
		assert.Nil(t, q)
	})
}

func TestDrivers(t *testing.T) {
	assert.Equal(t, []string{"mock", "test"}, Drivers())
}
