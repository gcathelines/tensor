package open_meteo

import (
	"errors"
	"testing"

	"github.com/gcathelines/tensor-energy-case/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestHourlyData_ToWeatherForecast(t *testing.T) {
	tests := []struct {
		name      string
		data      HourlyData
		expected  []types.WeatherForecast
		expectErr error
	}{
		{
			name: "failed, invalid data count",
			data: HourlyData{
				Time:          []string{"2024-09-06T00:00"},
				Temperature:   []float64{0.0},
				Precipitation: []float64{1.1, 2.2},
				WindSpeed:     []float64{},
				WindDirection: []float64{},
			},
			expectErr: errors.New("invalid data length, time 1, temp 1, precipitation 2, wind speed 0, wind direction 0"),
		},
		{
			name: "success",
			data: HourlyData{
				Time:          []string{"2024-09-06T00:00", "2024-09-06T01:00"},
				Temperature:   []float64{0.0, 1.0},
				Precipitation: []float64{1.2, 1.1},
				WindSpeed:     []float64{3.4, 1.2},
				WindDirection: []float64{5.6, 1.3},
			},
			expected: []types.WeatherForecast{
				{
					Time:          "2024-09-06T00:00",
					Temperature:   0.0,
					Precipitation: 1.2,
					WindSpeed:     3.4,
					WindDirection: 5.6,
				},
				{
					Time:          "2024-09-06T01:00",
					Temperature:   1.0,
					Precipitation: 1.1,
					WindSpeed:     1.2,
					WindDirection: 1.3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			forecasts, err := tt.data.ToWeatherForecasts()
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, forecasts); diff != "" {
				t.Fatalf("unexpected response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDailyData_HasPrecipitationToday(t *testing.T) {
	tests := []struct {
		name      string
		data      DailyData
		expected  bool
		expectErr error
	}{
		{
			name: "failed, invalid data count",
			data: DailyData{
				Time:             []string{"2024-09-06"},
				PrecipitationSum: []float64{1.1, 2.2},
			},
			expectErr: errors.New("invalid data length time 1, precipitation 2"),
		},
		{
			name: "success, has precipitations",
			data: DailyData{
				Time:             []string{"2024-09-06", "2024-09-07"},
				PrecipitationSum: []float64{1.1, 2.2},
			},
			expected: true,
		},
		{
			name: "success, has no precipitation",
			data: DailyData{
				Time:             []string{"2024-09-06", "2024-09-07"},
				PrecipitationSum: []float64{0, 1.1},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			precipitation, err := tt.data.HasPrecipitationToday()
			if tt.expectErr != nil {
				if err == nil || err.Error() != tt.expectErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if precipitation != tt.expected {
				t.Fatalf("expected %v got %v", tt.expected, precipitation)
			}
		})
	}
}
