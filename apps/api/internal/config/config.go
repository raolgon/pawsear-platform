package config

import (
	"errors"
	"os"
)

const (
	defaultHTTPAddr     = ":8080"
	defaultDatabasePath = "../../data/pawsear-local.db"
)

type Config struct {
	HTTPAddr     string
	DatabasePath string
	SeedDemoData bool
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:     envOrDefault("PAWSEAR_HTTP_ADDR", defaultHTTPAddr),
		DatabasePath: envOrDefault("PAWSEAR_DB_PATH", defaultDatabasePath),
		SeedDemoData: envOrDefault("PAWSEAR_SEED_DEMO", "false") == "true",
	}

	if cfg.HTTPAddr == "" {
		return Config{}, errors.New("PAWSEAR_HTTP_ADDR cannot be empty")
	}
	if cfg.DatabasePath == "" {
		return Config{}, errors.New("PAWSEAR_DB_PATH cannot be empty")
	}

	return cfg, nil
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
