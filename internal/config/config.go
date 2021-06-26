package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration structure
type Config struct {
	PostgresDSN string
}

// Load loads the application configuration from the environment variables or a .env file
func Load() *Config {
	godotenv.Load()

	return &Config{
		PostgresDSN: os.Getenv("IPR_POSTGRES_DSN"),
	}
}
