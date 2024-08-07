package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Run("invalid options", func(t *testing.T) {
		var options = map[string]any{
			"db": "invalid",
		}
		q, err := Open(options)
		assert.Error(t, err)
		assert.Nil(t, q)
	})

	t.Run("valid options", func(t *testing.T) {
		var options = map[string]any{
			"db":   testDB,
			"name": "test",
			"job":  new(testJob),
		}
		q, err := Open(options)
		assert.NoError(t, err)
		assert.NotNil(t, q)
	})
}
