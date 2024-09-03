package redis

import (
	"errors"
	queuecontract "github.com/gopi-frame/contract/queue"
	rediscontract "github.com/gopi-frame/contract/redis"
	"github.com/gopi-frame/exception"
)

// Config for redis queue
type Config struct {
	DB   rediscontract.Client
	Name string
	Job  queuecontract.Job
}

// NewConfig creates a new redis queue config
func NewConfig(db rediscontract.Client, name string, job queuecontract.Job) *Config {
	return &Config{
		DB:   db,
		Name: name,
		Job:  job,
	}
}

type Option func(cfg *Config) error

func (c *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

// Valid validates the config
func (c *Config) Valid() error {
	var errs []error
	if c.DB == nil {
		errs = append(errs, exception.New("db is required"))
	}
	if c.Name == "" {
		errs = append(errs, exception.New("name is required"))
	}
	if c.Job == nil {
		errs = append(errs, exception.New("job is required"))
	}
	return errors.Join(errs...)
}

// WithName sets the queue name. If name is empty, then name is not changed.
func WithName(name string) Option {
	return func(cfg *Config) error {
		if name == "" {
			return nil
		}
		cfg.Name = name
		return nil
	}
}

// WithDB sets the redis client. If db is nil, then db is not changed.
func WithDB(db rediscontract.Client) Option {
	return func(cfg *Config) error {
		if db == nil {
			return nil
		}
		cfg.DB = db
		return nil
	}
}

// WithJob sets the job type. If job is nil, then an exception is returned.
func WithJob(job queuecontract.Job) Option {
	return func(cfg *Config) error {
		if job == nil {
			return exception.NewEmptyArgumentException("job")
		}
		cfg.Job = job
		return nil
	}
}
