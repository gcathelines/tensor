package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gcathelines/tensor-energy-case/internal/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestUsecase_CreatePowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		testName  string
		name      string
		lat       float64
		long      float64
		expectErr error
	}{
		{
			testName: "success",
			name:     "My Cool Power Plant",
			lat:      1.1,
			long:     2.2,
		},
		{
			testName:  "failed, empty name",
			lat:       1.1,
			long:      2.2,
			expectErr: errors.New("name is required"),
		},
		{
			testName:  "failed, empty latitude",
			name:      "My Cool Power Plant",
			long:      2.2,
			expectErr: errors.New("latitude is required"),
		},
		{
			testName:  "failed, empty longitude",
			name:      "My Cool Power Plant",
			lat:       1.1,
			expectErr: errors.New("longitude is required"),
		},
		{
			testName:  "failed, invalid latitude",
			name:      "My Cool Power Plant",
			lat:       91.1,
			long:      2.2,
			expectErr: types.ErrInvalidLatitude,
		},
		{
			testName:  "failed, invalid longitude",
			name:      "My Cool Power Plant",
			lat:       11.1,
			long:      181.2,
			expectErr: types.ErrInvalidLongitude,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := testUsecase.CreatePowerPlant(ctx, tt.name, tt.lat, tt.long)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

}

func TestUsecase_UpdatePowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		testName  string
		name      string
		lat       float64
		long      float64
		id        int64
		expectErr error
	}{
		{
			testName: "success",
			name:     "My Cool Power Plant",
			lat:      1.1,
			long:     2.2,
			id:       1,
		},
		{
			testName:  "failed, invalid id",
			name:      "My Cool Power Plant",
			lat:       1.1,
			long:      2.2,
			id:        999,
			expectErr: errors.New("id not found"),
		},
		{
			testName:  "failed, empty id",
			name:      "My Cool Power Plant",
			lat:       1.1,
			long:      2.2,
			expectErr: errors.New("id is required"),
		},
		{
			testName:  "failed, invalid latitude",
			name:      "My Cool Power Plant",
			lat:       91.1,
			long:      2.2,
			id:        1,
			expectErr: types.ErrInvalidLatitude,
		},
		{
			testName:  "failed, invalid longitude",
			name:      "My Cool Power Plant",
			lat:       11.1,
			long:      181.2,
			id:        1,
			expectErr: types.ErrInvalidLongitude,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := testUsecase.UpdatePowerPlant(ctx, tt.id, &tt.name, &tt.lat, &tt.long)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUsecase_GetPowerPlant(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		testName     string
		id           int64
		forecastDays int
		expected     *types.PowerPlant
		expectErr    error
	}{
		{
			testName:     "success",
			id:           1,
			forecastDays: 7,
			expected: &types.PowerPlant{
				ID:        1,
				Name:      "My Cool Power Plant",
				Latitude:  22.11,
				Longitude: 33.11,
				Elevation: 0.6677740863787376,
				WeatherForecastProperties: types.WeatherForecastProperties{
					WeatherForecasts: []types.WeatherForecast{
						{
							Time:          "2024-09-06T00:00",
							Temperature:   1.1,
							Precipitation: 2.2,
							WindSpeed:     3.3,
							WindDirection: 4.4,
						},
						{
							Time:          "2024-09-06T01:00",
							Temperature:   11.1,
							Precipitation: 21.2,
							WindSpeed:     31.3,
							WindDirection: 41.4,
						},
					},
					HasPrecipitationToday: true,
				},
			},
		},
		{
			testName:  "failed, empty id",
			expectErr: errors.New("id is required"),
		},
		{
			testName:     "failed, invalid id",
			id:           999,
			forecastDays: 7,
			expectErr:    errors.New("id not found"),
		},
		{
			testName:     "failed, invalid forecast days",
			id:           1,
			forecastDays: 8,
			expectErr:    types.ErrInvalidForecastDay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			powerPlant, err := testUsecase.GetPowerPlant(ctx, tt.id, tt.forecastDays)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, powerPlant, cmpopts.IgnoreFields(types.PowerPlant{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUsecase_GetPowerPlants(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		testName     string
		lastID       int64
		count        int
		forecastDays int
		expected     []types.PowerPlant
		expectErr    error
	}{
		{
			testName:     "success",
			lastID:       0,
			count:        2,
			forecastDays: 7,
			expected: []types.PowerPlant{
				{
					ID:        1,
					Name:      "My Cool Power Plant 1",
					Latitude:  10.22,
					Longitude: 10.44,
					Elevation: 0.9789272030651343,
					WeatherForecastProperties: types.WeatherForecastProperties{
						WeatherForecasts: []types.WeatherForecast{
							{
								Time:          "2024-09-06T00:00",
								Temperature:   0.1,
								Precipitation: 0.2,
								WindSpeed:     0.3,
								WindDirection: 0.4,
							},
							{
								Time:          "2024-09-06T01:00",
								Temperature:   0.1,
								Precipitation: 0.2,
								WindSpeed:     0.3,
								WindDirection: 0.4,
							},
						},
						HasPrecipitationToday: true,
					},
				},
				{
					ID:        2,
					Name:      "My Cool Power Plant 2",
					Latitude:  20.22,
					Longitude: 20.44,
					Elevation: 0.9892367906066535,
					WeatherForecastProperties: types.WeatherForecastProperties{
						WeatherForecasts: []types.WeatherForecast{
							{
								Time:          "2024-09-06T00:00",
								Temperature:   1.1,
								Precipitation: 1.2,
								WindSpeed:     1.3,
								WindDirection: 1.4,
							},
							{
								Time:          "2024-09-06T01:00",
								Temperature:   10.1,
								Precipitation: 10.2,
								WindSpeed:     10.3,
								WindDirection: 10.4,
							},
						},
						HasPrecipitationToday: false,
					},
				},
			},
		},
		{
			testName:     "failed, invalid forecast days",
			lastID:       1,
			forecastDays: 8,
			expectErr:    types.ErrInvalidForecastDay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			powerPlants, err := testUsecase.GetPowerPlants(ctx, tt.lastID, tt.count, tt.forecastDays)
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, powerPlants, cmpopts.IgnoreFields(types.PowerPlant{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Fatalf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}
