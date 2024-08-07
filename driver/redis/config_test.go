package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig(testDB, "test", new(testJob))
	assert.Equal(t, "test", cfg.Name)
	assert.IsType(t, new(testJob), cfg.Job)
	assert.Equal(t, testDB, cfg.DB)
}

func TestConfig(t *testing.T) {
	t.Run("apply empty config", func(t *testing.T) {
		cfg := new(Config)
		err := cfg.Apply(
			WithName(""),
			WithDB(nil),
		)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "", cfg.Name)
		assert.Nil(t, cfg.DB)
	})

	t.Run("apply non-empty config", func(t *testing.T) {
		cfg := new(Config)
		err := cfg.Apply(
			WithName("test"),
			WithDB(testDB),
		)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "test", cfg.Name)
		assert.Equal(t, testDB, cfg.DB)
	})

	t.Run("apply empty job", func(t *testing.T) {
		cfg := new(Config)
		assert.Error(t, cfg.Apply(
			WithJob(nil),
		))
	})
}

func TestConfig_Valid(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := new(Config)
		err := cfg.Apply(
			WithName("test"),
			WithDB(testDB),
			WithJob(new(testJob)),
		)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Nil(t, cfg.Valid())
	})

	t.Run("invalid config", func(t *testing.T) {
		cfg := new(Config)
		assert.Error(t, cfg.Valid())
	})
}
