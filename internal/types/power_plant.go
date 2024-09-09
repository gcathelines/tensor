package types

import (
	"errors"
	"time"
)

var (
	ValidForecastLengths = map[int]struct{}{
		1:  {},
		3:  {},
		7:  {},
		14: {},
		16: {},
	}
	ErrInternal = errors.New("internal error")
)

type PowerPlant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Elevation float64   `json:"elevation"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	WeatherForecastProperties
}

type WeatherForecastProperties struct {
	HasPrecipitationToday bool              `json:"hasPrecipitationToday"`
	WeatherForecasts      []WeatherForecast `json:"weatherForecasts"`
}

type WeatherForecast struct {
	Time          string  `json:"time"`
	Temperature   float64 `json:"temperature"`
	Precipitation float64 `json:"precipitation"`
	WindSpeed     float64 `json:"windSpeed"`
	WindDirection float64 `json:"windDirection"`
}
