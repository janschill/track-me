package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	SentryDsn string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		DatabaseURL: os.Getenv("DB_PATH"),
		SentryDsn: os.Getenv("SENTRY_DSN"),
	}, nil
}
