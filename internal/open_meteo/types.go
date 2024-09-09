package open_meteo

import (
	"errors"
	"fmt"

	"github.com/gcathelines/tensor-energy-case/internal/types"
)

// ErrorResponse is generic response from the OpenMeteo API in case of an error.
type ErrorResponse struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

// WeatherForecast represents the response of Weather Forecast API
// Docs: https://open-meteo.com/en/docs
type WeatherForecast struct {
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	Elevation            float64    `json:"elevation"`
	GenerationtimeMs     float64    `json:"generationtime_ms"`
	UTCOffsetSeconds     int64      `json:"utc_offset_seconds"`
	Timezone             string     `json:"timezone"`
	TimezoneAbbreviation string     `json:"timezone_abbreviation"`
	Hourly               HourlyData `json:"hourly"`
	Daily                DailyData  `json:"daily"`
}

func (w WeatherForecast) ToProperties() (*types.WeatherForecastProperties, error) {
	forecasts, err := w.Hourly.ToWeatherForecasts()
	if err != nil {
		return nil, err
	}

	hasPrecipitationToday, err := w.Daily.HasPrecipitationToday()
	if err != nil {
		return nil, err
	}

	return &types.WeatherForecastProperties{
		WeatherForecasts:      forecasts,
		HasPrecipitationToday: hasPrecipitationToday,
	}, nil
}

type HourlyData struct {
	Time          []string  `json:"time"`
	Temperature   []float64 `json:"temperature_2m"`
	Precipitation []float64 `json:"precipitation"`
	WindSpeed     []float64 `json:"wind_speed_10m"`
	WindDirection []float64 `json:"wind_direction_10m"`
}

type DailyData struct {
	Time             []string  `json:"time"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
}

// ToWeatherForecasts converts the HourlyData to WeatherForecasts.
func (d HourlyData) ToWeatherForecasts() ([]types.WeatherForecast, error) {
	dataCount := len(d.Time)
	if len(d.Temperature) != dataCount ||
		len(d.Precipitation) != dataCount ||
		len(d.WindSpeed) != dataCount ||
		len(d.WindDirection) != dataCount {
		msg := fmt.Sprintf("invalid data length, time %d, temp %d, precipitation %d, wind speed %d, wind direction %d",
			dataCount, len(d.Temperature), len(d.Precipitation), len(d.WindSpeed), len(d.WindDirection))
		return nil, errors.New(msg)
	}

	forecasts := make([]types.WeatherForecast, 0, dataCount)
	for i := 0; i < dataCount; i++ {
		forecasts = append(forecasts, types.WeatherForecast{
			Time:          d.Time[i],
			Temperature:   d.Temperature[i],
			Precipitation: d.Precipitation[i],
			WindSpeed:     d.WindSpeed[i],
			WindDirection: d.WindDirection[i],
		})
	}
	return forecasts, nil
}

// HasPrecipitationToday returns true if there is precipitation today.
func (d DailyData) HasPrecipitationToday() (bool, error) {
	if len(d.Time) != len(d.PrecipitationSum) || len(d.Time) == 0 {
		msg := fmt.Sprintf("invalid data length time %d, precipitation %d",
			len(d.Time), len(d.PrecipitationSum))
		return false, errors.New(msg)
	}

	var hasPrecipitationToday bool
	// First daily data is always today
	if d.PrecipitationSum[0] > 0 {
		hasPrecipitationToday = true
	}

	return hasPrecipitationToday, nil
}

// Elevation represents the response of Elevation API
// Docs: https://open-meteo.com/en/docs/elevation-api
type Elevation struct {
	Elevation []float64 `json:"elevation"`
}
