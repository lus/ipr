package postgres

import (
	"context"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johejo/golang-migrate-extra/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

// postgresDriver represents the PostgreSQL data repository driver implementation
type postgresDriver struct {
	dsn      string
	pool     *pgxpool.Pool
	Machines *machineRepository
}

// New opens a new PostgreSQL connection pool and creates a new driver from it
func New(dsn string) (*postgresDriver, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &postgresDriver{
		dsn:      dsn,
		pool:     pool,
		Machines: &machineRepository{pool: pool},
	}, nil
}

// Migrate runs all pending migrations on the PostgreSQL database
func (driver *postgresDriver) Migrate() error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, driver.dsn)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Close closes the PostgreSQL connection pool and its driver
func (driver *postgresDriver) Close() {
	driver.pool.Close()
}
