package database

import (
	"context"
	"database/sql"

	"github.com/gcathelines/tensor-energy-case/internal/types"
)

// Database represents the database repository.
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database repository.
func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// CreatePowerPlant creates a new power plant in the database.
// This function returns the created power plant with the generated ID and version.
// Default version is 1.
func (d *Database) CreatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error) {
	query := `INSERT INTO power_plants (name, latitude, longitude) 
			VALUES ($1, $2, $3)
			RETURNING id, name, latitude, longitude, version, created_at`

	rows := d.db.QueryRowContext(ctx, query,
		powerPlant.Name,
		powerPlant.Latitude,
		powerPlant.Longitude,
	)

	var data types.PowerPlant
	err := rows.Scan(
		&data.ID,
		&data.Name,
		&data.Latitude,
		&data.Longitude,
		&data.Version,
		&data.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdatePowerPlant updates an existing power plant in the database.
// We use the version field to implement optimistic locking.
// This function returns the updated power plant with the new version.
func (d *Database) UpdatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error) {
	query := `UPDATE power_plants 
	SET name = $1, latitude = $2, longitude = $3, version = version + 1, updated_at = NOW()
	WHERE id = $4 and version = $5
	RETURNING id, name, latitude, longitude, version, created_at, updated_at`

	rows := d.db.QueryRowContext(ctx, query,
		powerPlant.Name,
		powerPlant.Latitude,
		powerPlant.Longitude,
		powerPlant.ID,
		powerPlant.Version,
	)

	var data types.PowerPlant
	err := rows.Scan(
		&data.ID,
		&data.Name,
		&data.Latitude,
		&data.Longitude,
		&data.Version,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetPowerPlant returns the power plant with the given ID.
func (d *Database) GetPowerPlant(ctx context.Context, id int64) (*types.PowerPlant, error) {
	query := `SELECT id, name, latitude, longitude, version, created_at, updated_at
	FROM power_plants WHERE id = $1`

	rows := d.db.QueryRowContext(ctx, query, id)

	var (
		data      types.PowerPlant
		updatedAt sql.NullTime
	)
	err := rows.Scan(
		&data.ID,
		&data.Name,
		&data.Latitude,
		&data.Longitude,
		&data.Version,
		&data.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	if updatedAt.Valid {
		data.UpdatedAt = updatedAt.Time
	}

	return &data, nil
}

// GetPowerPlants returns the power plants with the given last ID and count.
// The power plants are ordered by ID in ascending order.
func (d *Database) GetPowerPlants(ctx context.Context, lastID int64, count int) ([]types.PowerPlant, error) {
	query := `SELECT id, name, latitude, longitude, version, created_at, updated_at
	FROM power_plants WHERE id > $1
	ORDER BY id 
	FETCH FIRST $2 ROWS ONLY`

	powerPlants := []types.PowerPlant{}
	rows, err := d.db.QueryContext(ctx, query, lastID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			data      types.PowerPlant
			updatedAt sql.NullTime
		)
		err := rows.Scan(
			&data.ID,
			&data.Name,
			&data.Latitude,
			&data.Longitude,
			&data.Version,
			&data.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			data.UpdatedAt = updatedAt.Time
		}

		powerPlants = append(powerPlants, data)
	}

	return powerPlants, nil
}
