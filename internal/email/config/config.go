package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SmtpAddr string
	Email    string
	Password string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		SmtpAddr: os.Getenv("SMTP_ADDR"),
		Email:    os.Getenv("EMAIL"),
		Password: os.Getenv("PASSWORD"),
	}, nil
}
