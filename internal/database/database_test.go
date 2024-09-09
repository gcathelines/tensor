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

	closeFn := func(db *sql.DB) {
		db.Close()
	}

	return db, closeFn
}

func TestDatabase_CreatePowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tests := []struct {
		name      string
		payload   *types.PowerPlant
		expected  *types.PowerPlant
		expectErr error
	}{
		{
			name: "success",
			payload: &types.PowerPlant{
				Name:      "power plant 1",
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
			expected: &types.PowerPlant{
				Name:      "power plant 1",
				Latitude:  48.8566,
				Longitude: 2.3522,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powerPlant, err := testDB.CreatePowerPlant(ctx, tt.payload)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, powerPlant,
				cmpopts.IgnoreFields(types.PowerPlant{}, "ID", "CreatedAt")); diff != "" {
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
		expectErr   error
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
			},
			expected: &types.PowerPlant{
				Name:      "updated pp 1",
				Latitude:  1.1,
				Longitude: 2.2,
			},
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
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
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
		name      string
		ID        int64
		expected  *types.PowerPlant
		expectErr error
	}{
		{
			name: "success",
			ID:   powerPlant.ID,
			expected: &types.PowerPlant{
				ID:        powerPlant.ID,
				Name:      "Solar Power Plant",
				Latitude:  50.8503,
				Longitude: 4.3517,
			},
		},
		{
			name:      "not found",
			ID:        0,
			expectErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powerPlant, err := testDB.GetPowerPlant(ctx, tt.ID)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
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

func TestDatabase_GetPowerPlantForUpdate(t *testing.T) {
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
		name      string
		ID        int64
		expected  *types.PowerPlant
		expectErr error
	}{
		{
			name: "success",
			ID:   powerPlant.ID,
			expected: &types.PowerPlant{
				ID:        powerPlant.ID,
				Name:      "Solar Power Plant",
				Latitude:  50.8503,
				Longitude: 4.3517,
			},
		},
		{
			name:      "not found",
			ID:        0,
			expectErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powerPlant, err := testDB.GetPowerPlantForUpdate(ctx, tt.ID)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
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

	tests := []struct {
		name      string
		lastID    int64
		count     int
		expected  []types.PowerPlant
		expectErr error
	}{
		{
			name:   "page 1, count 3",
			lastID: 0,
			count:  3,
			expected: []types.PowerPlant{
				{
					ID:        1,
					Name:      "Solar Power Plant",
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				{
					ID:        2,
					Name:      "Wind Power Plant",
					Latitude:  34.0522,
					Longitude: -118.2437,
				},
				{
					ID:        3,
					Name:      "Hydro Power Plant",
					Latitude:  37.7749,
					Longitude: -122.4194,
				},
			},
		},
		{
			name:   "page 2, count 2",
			lastID: 3,
			count:  2,
			expected: []types.PowerPlant{
				{
					ID:        4,
					Name:      "Solar 2 Power Plant",
					Latitude:  40.7128,
					Longitude: -74.0060,
				},
				{
					ID:        5,
					Name:      "Wind 2 Power Plant",
					Latitude:  34.0522,
					Longitude: -118.2437,
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
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
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
