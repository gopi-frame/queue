package redis

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Run("invalid options", func(t *testing.T) {
		options := map[string]any{
			"db": "invalid",
		}
		q, err := Open(options)
		assert.Error(t, err)
		assert.Nil(t, q)
	})

	t.Run("valid options", func(t *testing.T) {
		options := map[string]any{
			"db":   testDB,
			"name": "test",
			"job":  new(testJob),
		}
		q, err := Open(options)
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			assert.Equal(t, "test", q.Name())
			assert.Equal(t, testDB, q.(*Queue).client)
			assert.Equal(t, reflect.TypeFor[testJob](), q.(*Queue).job)
		}
	})
}
