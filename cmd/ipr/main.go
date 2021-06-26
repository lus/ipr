package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lus/ipr/internal/config"
	"github.com/lus/ipr/internal/database/postgres"
	"github.com/lus/ipr/internal/server"
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

	// Run the web server
	go func() {
		settings := &server.Settings{
			Address:   cfg.Address,
			AuthToken: cfg.AuthToken,
		}
		repositories := &server.Repositories{
			MachineRepository: driver.Machines,
		}

		log.Info().Str("address", cfg.Address).Msg("Starting the web server...")
		err := server.RunBlocking(settings, repositories)
		log.Info().Msg("Shutting down the web server...")
		if err != nil {
			log.Error().Err(err).Msg("Could not shut down the web server")
		}
	}()

	// Wait for the program to exit
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	<-channel
}
