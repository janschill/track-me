package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	expectedDBPath := "test_db_path"
	expectedSentryDsn := "test_sentry_dsn"
	os.Setenv("DB_PATH", expectedDBPath)
	os.Setenv("SENTRY_DSN", expectedSentryDsn)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, wantErr %v", err, false)
	}

	if cfg.DatabaseURL != expectedDBPath {
		t.Errorf("LoadConfig().DatabaseURL = %v, want %v", cfg.DatabaseURL, expectedDBPath)
	}

	if cfg.SentryDsn != expectedSentryDsn {
		t.Errorf("LoadConfig().SentryDsn = %v, want %v", cfg.SentryDsn, expectedSentryDsn)
	}

	os.Unsetenv("DB_PATH")
	os.Unsetenv("SENTRY_DSN")
}
