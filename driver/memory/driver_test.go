package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Run("invalid options", func(t *testing.T) {
		options := map[string]any{
			"name": 123,
		}
		assert.Panics(t, func() {
			_, _ = Open(options)
		})
	})

	t.Run("valid options", func(t *testing.T) {
		options := map[string]any{
			"name": "test",
		}
		q, err := Open(options)
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.Equal(t, "test", q.Name())
		}
	})

	t.Run("empty name", func(t *testing.T) {
		options := map[string]any{
			"name": "",
		}
		q, err := Open(options)
		assert.Error(t, err)
		assert.Nil(t, q)
	})

	t.Run("undefined index name", func(t *testing.T) {
		options := map[string]any{}
		q, err := Open(options)
		assert.Error(t, err)
		assert.Nil(t, q)
	})
}
