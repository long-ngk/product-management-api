package config

import (
	"os"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load reads configuration from environment variables.
// DATABASE_URL is required; PORT defaults to ":8080" if not set.
func Load() (*Config, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, &MissingEnvError{Key: "DATABASE_URL"}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	return &Config{
		DatabaseURL: dsn,
		Port:        port,
	}, nil
}

// MissingEnvError is returned when a required environment variable is absent.
type MissingEnvError struct {
	Key string
}

func (e *MissingEnvError) Error() string {
	return e.Key + " environment variable is required"
}
