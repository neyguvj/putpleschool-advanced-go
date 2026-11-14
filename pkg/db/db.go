package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxOpenConns    = 100
	maxIdleConns    = 10
	connMaxLifetime = time.Hour
)

func New(cfg *Config) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.New(postgres.Config{DSN: cfg.DSN()}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return gdb, nil
}
