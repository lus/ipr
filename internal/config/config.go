package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration structure
type Config struct {
	PostgresDSN string
	Address     string
	AuthToken   string
}

// Load loads the application configuration from the environment variables or a .env file
func Load() *Config {
	godotenv.Load()

	return &Config{
		PostgresDSN: envString("IPR_POSTGRES_DSN", ""),
		Address:     envString("IPR_ADDRESS", ":8080"),
		AuthToken:   envString("IPR_AUTH_TOKEN", "foobar"),
	}
}

func envString(key, fallback string) string {
	value, set := os.LookupEnv(key)
	if !set {
		return fallback
	}
	return value
}
