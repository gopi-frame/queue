package database

import (
	"errors"
	"github.com/gopi-frame/contract/queue"
	"github.com/gopi-frame/exception"
	"gorm.io/gorm"
)

// DefaultJobTable is the default Table name for jobs.
const DefaultJobTable = "jobs"

// Config is the configuration for database queue.
type Config struct {
	DB    *gorm.DB
	Name  string
	Table string
	Job   queue.Job
}

// Valid validates the config. It returns an exception if the config is invalid.
func (c *Config) Valid() error {
	var errs []error
	if c.DB == nil {
		errs = append(errs, exception.New("db is nil"))
	}
	if c.Name == "" {
		errs = append(errs, exception.New("name is empty"))
	}
	if c.Job == nil {
		errs = append(errs, exception.New("job is nil"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// Option is the option for database queue.
type Option func(cfg *Config) error

// NewConfig returns a new config.
func NewConfig(db *gorm.DB, name string, job queue.Job) *Config {
	return &Config{
		DB:    db,
		Name:  name,
		Table: DefaultJobTable,
		Job:   job,
	}
}

// Apply applies options to config.
func (c *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}
	return nil
}

// WithTable sets the Table name.
func WithTable(table string) Option {
	return func(cfg *Config) error {
		if table == "" {
			return nil
		}
		cfg.Table = table
		return nil
	}
}

// WithName sets the queue name.
func WithName(name string) Option {
	return func(cfg *Config) error {
		if name == "" {
			return nil
		}
		cfg.Name = name
		return nil
	}
}

// WithDB sets the database.
func WithDB(db *gorm.DB) Option {
	return func(cfg *Config) error {
		if db == nil {
			return nil
		}
		cfg.DB = db
		return nil
	}
}

// WithJob sets the job type.
func WithJob(job queue.Job) Option {
	return func(cfg *Config) error {
		if job == nil {
			return exception.NewEmptyArgumentException("job")
		}
		cfg.Job = job
		return nil
	}
}
