package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gcathelines/tensor-energy-case/internal/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
)

var (
	testDB *Database
)

func TestMain(m *testing.M) {
	db, close := SetupDB()

	testDB = NewDatabase(db)
	code := m.Run()

	close(db)
	os.Exit(code)
}

func SetupDB() (*sql.DB, func(*sql.DB)) {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/testdb?sslmode=disable")
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	schemaUp, err := os.ReadFile("./../../migrations/schema.up.sql")
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	_, err = db.Exec(string(schemaUp))
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	closeFn := func(db *sql.DB) {
		schemaDown, err := os.ReadFile("./../../migrations/schema.down.sql")
		if err != nil {
			db.Close()
			log.Fatal(err)
		}

		_, err = db.Exec(string(schemaDown))
		if err != nil {
			db.Close()
			log.Fatal(err)
		}

		db.Close()
	}

	return db, closeFn
}

func TestDatabase_CreatePowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		payload  *types.PowerPlant
		expected *types.PowerPlant
		err      error
	}{
		{
			name: "success",
			payload: &types.PowerPlant{
				Name:      "power plant 1",
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
			expected: &types.PowerPlant{
				ID:        1,
				Name:      "power plant 1",
				Latitude:  48.8566,
				Longitude: 2.3522,
				Version:   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powerPlant, err := testDB.CreatePowerPlant(ctx, tt.payload)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, powerPlant,
				cmpopts.IgnoreFields(types.PowerPlant{}, "CreatedAt")); diff != "" {
				t.Fatalf("unexpected power plant (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDatabase_UpdatePowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tests := []struct {
		name        string
		createdData *types.PowerPlant
		payload     *types.PowerPlant
		expected    *types.PowerPlant
		err         error
	}{
		{
			name: "success",
			createdData: &types.PowerPlant{
				Name:      "power plant 1",
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
			payload: &types.PowerPlant{
				Name:      "updated pp 1",
				Latitude:  1.1,
				Longitude: 2.2,
				Version:   1,
			},
			expected: &types.PowerPlant{
				Name:      "updated pp 1",
				Latitude:  1.1,
				Longitude: 2.2,
				Version:   2,
			},
		},
		{
			name: "failed, version mismatch",
			createdData: &types.PowerPlant{
				Name:      "power plant 2",
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
			payload: &types.PowerPlant{
				Name:      "updated pp 2",
				Latitude:  1.1,
				Longitude: 2.2,
				Version:   2,
			},
			err: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdPowerPlant, err := testDB.CreatePowerPlant(ctx, tt.createdData)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			tt.payload.ID = createdPowerPlant.ID

			powerPlant, err := testDB.UpdatePowerPlant(ctx, tt.payload)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if powerPlant.UpdatedAt.Equal(time.Time{}) {
				t.Fatalf("expected updated at to be set, got: %v", powerPlant.UpdatedAt)
			}

			if diff := cmp.Diff(tt.expected, powerPlant,
				cmpopts.IgnoreFields(types.PowerPlant{}, "ID", "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("unexpected power plant (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDatabase_GetPowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	powerPlant, err := testDB.CreatePowerPlant(ctx, &types.PowerPlant{
		Name:      "Solar Power Plant",
		Latitude:  50.8503,
		Longitude: 4.3517,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		ID       int64
		expected *types.PowerPlant
		err      error
	}{
		{
			name: "success",
			ID:   powerPlant.ID,
			expected: &types.PowerPlant{
				ID:        powerPlant.ID,
				Name:      "Solar Power Plant",
				Latitude:  50.8503,
				Longitude: 4.3517,
				Version:   1,
			},
		},
		{
			name: "not found",
			ID:   0,
			err:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powerPlant, err := testDB.GetPowerPlant(ctx, tt.ID)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, powerPlant,
				cmpopts.IgnoreFields(types.PowerPlant{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("unexpected power plant (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDatabase_GetPowerPlants(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := testDB.db.Exec("TRUNCATE TABLE power_plants")
	if err != nil {
		t.Fatal(err)
	}

	seedUp, err := os.ReadFile("./../../migrations/seed.up.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = testDB.db.Exec(string(seedUp))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		lastID   int64
		count    int
		expected []types.PowerPlant
		err      error
	}{
		{
			name:   "page 1, count 3",
			lastID: 100000000,
			count:  3,
			expected: []types.PowerPlant{
				{
					ID:        100000001,
					Name:      "Solar Power Plant",
					Latitude:  40.7128,
					Longitude: -74.0060,
					Version:   1,
				},
				{
					ID:        100000002,
					Name:      "Wind Power Plant",
					Latitude:  34.0522,
					Longitude: -118.2437,
					Version:   1,
				},
				{
					ID:        100000003,
					Name:      "Hydro Power Plant",
					Latitude:  37.7749,
					Longitude: -122.4194,
					Version:   1,
				},
			},
		},
		{
			name:   "page 2, count 2",
			lastID: 100000003,
			count:  2,
			expected: []types.PowerPlant{
				{
					ID:        100000004,
					Name:      "Solar 2 Power Plant",
					Latitude:  40.7128,
					Longitude: -74.0060,
					Version:   1,
				},
				{
					ID:        100000005,
					Name:      "Wind 2 Power Plant",
					Latitude:  34.0522,
					Longitude: -118.2437,
					Version:   1,
				},
			},
		},
		{
			name:   "page 4, count 3",
			lastID: 100000009,
			count:  3,
			expected: []types.PowerPlant{
				{
					ID:        100000010,
					Name:      "Last Power Plant",
					Latitude:  40.7128,
					Longitude: -74.0060,
					Version:   1,
				},
			},
		},
		{
			name:     "page 2, count 10",
			lastID:   100000010,
			count:    10,
			expected: []types.PowerPlant{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powerPlant, err := testDB.GetPowerPlants(ctx, tt.lastID, tt.count)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, powerPlant,
				cmpopts.IgnoreFields(types.PowerPlant{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("unexpected power plant (-want +got):\n%s", diff)
			}
		})
	}

	seedDown, err := os.ReadFile("./../../migrations/seed.down.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = testDB.db.Exec(string(seedDown))
	if err != nil {
		log.Fatal(err)
	}
}
