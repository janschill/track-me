package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	SentryDsn string
	GarminIPCInbound string
	GarminDeviceIMEI string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		DatabaseURL: os.Getenv("DB_PATH"),
		SentryDsn: os.Getenv("SENTRY_DSN"),
		GarminIPCInbound: os.Getenv("GARMIN_IPC_INBOUND"),
		GarminDeviceIMEI: os.Getenv("GARMIN_DEVICE_IMEI"),
	}, nil
}
