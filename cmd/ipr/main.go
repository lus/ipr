package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lus/ipr/internal/config"
	"github.com/lus/ipr/internal/database/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set up the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load the application configuration
	log.Info().Msg("Loading the application configuration...")
	cfg := config.Load()

	// Initialize the PostgreSQL repository driver
	log.Info().Str("dsn", cfg.PostgresDSN).Msg("Opening the PostgreSQL connection...")
	driver, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatal().Err(err).Str("dsn", cfg.PostgresDSN).Msg("Could not open the PostgreSQL connection")
	}
	defer func() {
		log.Info().Msg("Closing the PostgreSQL connection...")
		driver.Close()
	}()

	// Run migrations on the PostgreSQL driver
	log.Info().Msg("Running SQL migrations on the PostgreSQL driver...")
	if err := driver.Migrate(); err != nil {
		log.Fatal().Err(err).Msg("Could not run SQL migrations on the PostgreSQL driver")
	}

	// Wait for the program to exit
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	<-channel
}
