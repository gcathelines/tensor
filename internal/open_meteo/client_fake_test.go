package open_meteo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ServeFakeOpenMeteo serves a fake OpenMeteo API for testing purposes.
func ServeFakeOpenMeteo(t *testing.T, ctx context.Context) string {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		lats := r.URL.Query()["latitude"]
		longs := r.URL.Query()["longitude"]

		// For testing purpose, we set if latitude and longitude is 200.1
		// we return a bad request response.
		if len(lats) == 1 && lats[0] == "200.1" &&
			len(longs) == 1 && longs[0] == "200.1" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBadRequest)
			return
		}

		switch r.URL.Path {
		case "/v1/forecast":
			w.WriteHeader(http.StatusOK)

			if len(lats) > 1 && len(longs) > 1 {
				w.Write(responseForecasts)
			} else {
				w.Write(responseForecast)
			}
		case "/v1/elevation":
			w.WriteHeader(http.StatusOK)
			w.Write(responseElevation)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write(responseNotFound)
		}
	}))
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	return srv.URL
}

var (
	responseNotFound   = []byte(`{"error":true,"reason":"Not Found"}`)
	responseBadRequest = []byte(`{"error":true,"reason":"Parameter 'latitude' and 'longitude' must have the same number of elements"}`)
	responseElevation  = []byte(`{"elevation":[38.01]}`)
	responseForecast   = []byte(`
		{
			"latitude": 52.52,
			"longitude": 13.41,
			"generationtime_ms": 0.15592575073242188,
			"utc_offset_seconds": 0,
			"timezone": "GMT",
			"timezone_abbreviation": "GMT",
			"elevation": 38,
			"hourly_units": {
				"time": "iso8601",
				"temperature_2m": "°C",
				"precipitation": "mm",
				"wind_speed_10m": "km/h",
				"wind_direction_10m": "°"
			},
			"hourly": {
				"time": ["2024-09-06T00:00","2024-09-06T01:00","2024-09-06T02:00"],
				"temperature_2m": [22.1,21.2,20.5],
				"precipitation": [0.1,0.2,0.3],
				"wind_speed_10m": [11.9,12.4,12.8],
				"wind_direction_10m": [85,80,80]
			},
			"daily_units": {
				"time": "iso8601",
				"precipitation_sum": "mm"
			},
			"daily": {
				"time": [
					"2024-09-07",
					"2024-09-08",
					"2024-09-09"
				],
				"precipitation_sum": [
					0,
					0,
					8.8
				]
			}
		}
	`)
	responseForecasts = []byte(`[
		{
			"latitude": 52.52,
			"longitude": 13.419998,
			"generationtime_ms": 0.20694732666015625,
			"utc_offset_seconds": 0,
			"timezone": "GMT",
			"timezone_abbreviation": "GMT",
			"elevation": 38,
			"hourly_units": {
			"time": "iso8601",
			"temperature_2m": "°C",
			"precipitation": "mm",
			"wind_speed_10m": "km/h",
			"wind_direction_10m": "°"
			},
			"hourly": {
				"time": [
					"2024-09-07T00:00",
					"2024-09-07T01:00",
					"2024-09-07T02:00"
				],
				"temperature_2m": [
					18.4,
					17.8,
					17.3
				],
				"precipitation": [
					0.2,
					0.5,
					2.1
				],
				"wind_speed_10m": [
					5.9,
					6.6,
					7.1
				],
				"wind_direction_10m": [
					104,
					99,
					120
				]
			},
			"daily": {
				"time": [
					"2024-09-07",
					"2024-09-08",
					"2024-09-09"
				],
				"precipitation_sum": [
					0.1,
					0,
					2.1
				]
			}
		},
		{
			"latitude": 14.125,
			"longitude": 15.125,
			"generationtime_ms": 0.5570650100708008,
			"utc_offset_seconds": 0,
			"timezone": "GMT",
			"timezone_abbreviation": "GMT",
			"elevation": 333,
			"location_id": 1,
			"hourly_units": {
			"time": "iso8601",
			"temperature_2m": "°C",
			"precipitation": "mm",
			"wind_speed_10m": "km/h",
			"wind_direction_10m": "°"
			},
			"hourly": {
				"time": [
					"2024-09-07T00:00",
					"2024-09-07T01:00",
					"2024-09-07T02:00"
				],
				"precipitation": [
					2.7,
					2.6,
					0.5
				],
				"temperature_2m": [
					25.9,
					25.5,
					25.2
				],
				"wind_speed_10m": [
					9,
					8.4,
					6.9
				],
				"wind_direction_10m": [
					157,
					155,
					152
				]
			},
			"daily": {
				"time": [
					"2024-09-07",
					"2024-09-08",
					"2024-09-09"
				],
				"precipitation_sum": [
					0,
					0,
					0
				]
			}
		}
	]`)
)
