package types

import "time"

var (
	ValidForecastLengths = map[int]struct{}{
		1:  {},
		3:  {},
		7:  {},
		14: {},
		16: {},
	}
)

type PowerPlant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Elevation float64   `json:"elevation"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	WeatherForecastProperties
}

type WeatherForecastProperties struct {
	HasPrecipitationToday bool              `json:"has_precipitation_today"`
	WeatherForecasts      []WeatherForecast `json:"weather_forecasts"`
}

type WeatherForecast struct {
	Time          string  `json:"time"`
	Temperature   float64 `json:"temperature"`
	Precipitation float64 `json:"precipitation"`
	WindSpeed     float64 `json:"wind_speed"`
	WindDirection float64 `json:"wind_direction"`
}
