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
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	WeatherForecastProperties
}

type WeatherForecastProperties struct {
	HasPrecipitationToday bool              `json:"has_precipitation_today"`
	WeatherForecasts      []WeatherForecast `json:"weather_forecasts"`
}

type WeatherForecast struct {
	Time             string  `json:"time"`
	Temperature2m    float64 `json:"temperature_2m"`
	Precipitation    float64 `json:"precipitation"`
	WindSpeed10m     float64 `json:"wind_speed_10m"`
	WindDirection10m float64 `json:"wind_direction_10m"`
}
