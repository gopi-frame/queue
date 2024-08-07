package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig(testDB, "test", new(testJob))
	assert.Equal(t, "jobs", cfg.Table)
	assert.Equal(t, "test", cfg.Name)
	assert.IsType(t, new(testJob), cfg.Job)
	assert.Equal(t, testDB, cfg.DB)
}

func TestConfig_Apply(t *testing.T) {
	t.Run("apply empty config", func(t *testing.T) {
		cfg := new(Config)
		err := cfg.Apply(
			WithTable(""),
			WithName(""),
			WithDB(nil),
		)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "", cfg.Table)
		assert.Equal(t, "", cfg.Name)
		assert.Nil(t, cfg.DB)

		assert.Panics(t, func() {
			err := cfg.Apply(WithJob(nil))
			if err != nil {
				panic(err)
			}
		})
	})

	t.Run("apply non-empty config", func(t *testing.T) {
		cfg := new(Config)
		err := cfg.Apply(
			WithTable("jobs2"),
			WithName("test"),
			WithDB(testDB),
			WithJob(new(testJob)),
		)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, testDB, cfg.DB)
		assert.Equal(t, "jobs2", cfg.Table)
		assert.Equal(t, "test", cfg.Name)
		assert.IsType(t, new(testJob), cfg.Job)
	})
}
