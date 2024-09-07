package open_meteo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gcathelines/tensor-energy-case/config"
	"github.com/gcathelines/tensor-energy-case/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestOpenMeteoClient_GetWeatherForecast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fakeURL := ServeFakeOpenMeteo(t, ctx)
	cl := NewOpenMeteoClient(config.OpenMeteoConfig{
		APIURL:  fakeURL,
		Timeout: 5 * time.Second,
	})

	tests := []struct {
		name     string
		lat      float64
		long     float64
		err      error
		expected *types.WeatherForecastProperties
	}{
		{
			name:     "failed, bad request",
			lat:      200.1,
			long:     200.1,
			err:      errors.New("unexpected status code: 400, reason: Parameter 'latitude' and 'longitude' must have the same number of elements"),
			expected: nil,
		},
		{
			name: "success",
			lat:  52.52,
			long: 13.41,
			err:  nil,
			expected: &types.WeatherForecastProperties{
				HasPrecipitationToday: false,
				WeatherForecasts: []types.WeatherForecast{
					{
						Time:             "2024-09-06T00:00",
						Temperature2m:    22.1,
						Precipitation:    0.1,
						WindSpeed10m:     11.9,
						WindDirection10m: 85,
					},
					{
						Time:             "2024-09-06T01:00",
						Temperature2m:    21.2,
						Precipitation:    0.2,
						WindSpeed10m:     12.4,
						WindDirection10m: 80,
					},
					{
						Time:             "2024-09-06T02:00",
						Temperature2m:    20.5,
						Precipitation:    0.3,
						WindSpeed10m:     12.8,
						WindDirection10m: 80,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := cl.GetWeatherForecast(ctx, tt.lat, tt.long, 7)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, resp); diff != "" {
				t.Fatalf("unexpected response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOpenMeteoClient_GetWeatherForecasts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fakeURL := ServeFakeOpenMeteo(t, ctx)
	cl := NewOpenMeteoClient(config.OpenMeteoConfig{
		APIURL:  fakeURL,
		Timeout: 5 * time.Second,
	})

	tests := []struct {
		name     string
		lats     []float64
		longs    []float64
		err      error
		expected []types.WeatherForecastProperties
	}{
		{
			name:     "failed, bad request",
			lats:     []float64{200.1},
			longs:    []float64{200.1},
			err:      errors.New("unexpected status code: 400, reason: Parameter 'latitude' and 'longitude' must have the same number of elements"),
			expected: nil,
		},
		{
			name:  "success",
			lats:  []float64{52.52, 71.55},
			longs: []float64{13.419998, 62.01},
			err:   nil,
			expected: []types.WeatherForecastProperties{
				{
					HasPrecipitationToday: true,
					WeatherForecasts: []types.WeatherForecast{
						{
							Time:             "2024-09-07T00:00",
							Temperature2m:    18.4,
							Precipitation:    0.2,
							WindSpeed10m:     5.9,
							WindDirection10m: 104,
						},
						{
							Time:             "2024-09-07T01:00",
							Temperature2m:    17.8,
							Precipitation:    0.5,
							WindSpeed10m:     6.6,
							WindDirection10m: 99,
						},
						{
							Time:             "2024-09-07T02:00",
							Temperature2m:    17.3,
							Precipitation:    2.1,
							WindSpeed10m:     7.1,
							WindDirection10m: 120,
						},
					},
				},
				{
					HasPrecipitationToday: false,
					WeatherForecasts: []types.WeatherForecast{
						{
							Time:             "2024-09-07T00:00",
							Temperature2m:    25.9,
							Precipitation:    2.7,
							WindSpeed10m:     9,
							WindDirection10m: 157,
						},
						{
							Time:             "2024-09-07T01:00",
							Temperature2m:    25.5,
							Precipitation:    2.6,
							WindSpeed10m:     8.4,
							WindDirection10m: 155,
						},
						{
							Time:             "2024-09-07T02:00",
							Temperature2m:    25.2,
							Precipitation:    0.5,
							WindSpeed10m:     6.9,
							WindDirection10m: 152,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := cl.GetWeatherForecasts(ctx, tt.lats, tt.longs, 7)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, resp); diff != "" {
				t.Fatalf("unexpected response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOpenMeteoClient_GetElevation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fakeURL := ServeFakeOpenMeteo(t, ctx)
	cl := NewOpenMeteoClient(config.OpenMeteoConfig{
		APIURL:  fakeURL,
		Timeout: 5 * time.Second,
	})

	tests := []struct {
		name     string
		lat      float64
		long     float64
		err      error
		expected float64
	}{
		{
			name:     "failed, bad request",
			lat:      200.1,
			long:     200.1,
			err:      errors.New("unexpected status code: 400, reason: Parameter 'latitude' and 'longitude' must have the same number of elements"),
			expected: 0,
		},
		{
			name:     "success",
			lat:      47.36865,
			long:     8.539183,
			err:      nil,
			expected: 38.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := cl.GetElevation(ctx, tt.lat, tt.long)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, resp); diff != "" {
				t.Fatalf("unexpected response (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOpenMeteoClient_doRequest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fakeURL := ServeFakeOpenMeteo(t, ctx)
	cl := NewOpenMeteoClient(config.OpenMeteoConfig{
		APIURL:  fakeURL,
		Timeout: 5 * time.Second,
	})

	tests := []struct {
		name     string
		path     string
		err      error
		value    any
		expected any
	}{
		{
			name: "failed, not found",
			path: "/v1/unknown",
			err:  errors.New("unexpected status code: 404, reason: Not Found"),
		},
		{
			name:  "success",
			path:  "/v1/elevation",
			err:   nil,
			value: map[string]any{},
			expected: map[string]any{
				"elevation": []any{38.01},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cl.doRequest(tt.path, nil, "GET", nil, &tt.value)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error: %v, got: %v", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.expected, tt.value); diff != "" {
				t.Fatalf("unexpected response (-want +got):\n%s", diff)
			}
		})
	}
}
