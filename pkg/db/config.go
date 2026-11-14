package db

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST" validate:"required,hostname|ip"`
	Port     string `env:"POSTGRES_PORT" validate:"required,numeric"`
	User     string `env:"POSTGRES_USER" validate:"required"`
	Password string `env:"POSTGRES_PASSWORD" validate:"required"`
	DBName   string `env:"POSTGRES_DB" validate:"required"`
	SSLMode  string `env:"POSTGRES_SSLMODE" validate:"omitempty,oneof=disable allow prefer require verify-ca verify-full"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // ignore missing .env

	cfg := &Config{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		DBName:   getEnv("POSTGRES_DB", "advancedgo"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid db config: %w", err)
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
