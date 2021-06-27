package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lus/ipr/internal/shared"
)

// machineRepository represents the PostgreSQL machine repository implementation
type machineRepository struct {
	pool *pgxpool.Pool
}

// All looks up all stored machines
func (repository *machineRepository) All() ([]*shared.Machine, error) {
	query := "SELECT * FROM machines"

	rows, err := repository.pool.Query(context.Background(), query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*shared.Machine{}, nil
		}
		return nil, err
	}

	var machines []*shared.Machine
	for rows.Next() {
		machine, err := rowToMachine(rows)
		if err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

// Lookup looks up a stored machine with a specific name
func (repository *machineRepository) Lookup(name string) (*shared.Machine, error) {
	query := "SELECT * FROM machines WHERE name = $1"

	machine, err := rowToMachine(repository.pool.QueryRow(context.Background(), query, name))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return machine, nil
}

// Upsert creates or replaces a machine, depending on whether or not its name is taken
func (repository *machineRepository) Upsert(machine *shared.Machine) error {
	query := `
		INSERT INTO machines (name, token, address, updated)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE
			SET token = excluded.token,
				address = excluded.address,
				updated = excluded.updated
	`

	_, err := repository.pool.Exec(context.Background(), query, machine.Name, machine.Token, machine.Address, machine.Updated)
	return err
}

// Delete deletes a machine by its name
func (repository *machineRepository) Delete(name string) error {
	query := "DELETE FROM machines WHERE name = $1"

	_, err := repository.pool.Exec(context.Background(), query, name)
	return err
}

func rowToMachine(row pgx.Row) (*shared.Machine, error) {
	machine := new(shared.Machine)

	if err := row.Scan(&machine.Name, &machine.Token, &machine.Address, &machine.Updated); err != nil {
		return nil, err
	}

	return machine, nil
}
