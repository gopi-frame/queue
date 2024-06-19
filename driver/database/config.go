package database

import "gorm.io/gorm"

type Config struct {
	db    *gorm.DB
	name  string
	table string
}

type Option func(cfg *Config) error

func NewConfig(db *gorm.DB, name string) *Config {
	return &Config{
		db:    db,
		name:  name,
		table: DefaultJobTable,
	}
}

func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}

func Table(table string) Option {
	return func(cfg *Config) error {
		if table == "" {
			return nil
		}
		cfg.table = table
		return nil
	}
}

func (cfg *Config) FormatOpts() map[string]any {
	return map[string]any{
		OptKeyConnection: cfg.db,
		OptKeyName:       cfg.name,
		OptKeyTable:      cfg.table,
	}
}

func ParseCfg(options map[string]any) (*Config, error) {
	cfg := new(Config)
	if db, ok := options[OptKeyConnection]; ok {
		cfg.db = db.(*gorm.DB)
	} else {
		return nil, ErrMissingOptionConnection
	}
	if name, ok := options[OptKeyName]; ok {
		cfg.name = name.(string)
	} else {
		return nil, ErrMissingOptionName
	}
	if table, ok := options[OptKeyTable]; ok {
		cfg.table = table.(string)
	} else {
		cfg.table = DefaultJobTable
	}
	return cfg, nil
}
