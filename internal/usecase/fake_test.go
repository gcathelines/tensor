package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gcathelines/tensor-energy-case/internal/types"
)

var (
	testUsecase *Usecase
)

func TestMain(m *testing.M) {
	testUsecase = NewUsecase(&fakeWeatherAPI{}, &fakeDB{})
	code := m.Run()
	os.Exit(code)
}

type fakeWeatherAPI struct{}

func (f *fakeWeatherAPI) GetWeatherForecast(ctx context.Context, latitudes float64, longitudes float64, forecastDays int) (*types.WeatherForecastProperties, error) {
	return &types.WeatherForecastProperties{
		WeatherForecasts: []types.WeatherForecast{
			{
				Time:             "2024-09-06T00:00",
				Temperature2m:    1.1,
				Precipitation:    2.2,
				WindSpeed10m:     3.3,
				WindDirection10m: 4.4,
			},
			{
				Time:             "2024-09-06T01:00",
				Temperature2m:    11.1,
				Precipitation:    21.2,
				WindSpeed10m:     31.3,
				WindDirection10m: 41.4,
			},
		},
		HasPrecipitationToday: true,
	}, nil
}

func (f *fakeWeatherAPI) GetWeatherForecasts(ctx context.Context, latitudes []float64, longitudes []float64, forecastDays int) ([]types.WeatherForecastProperties, error) {
	count := len(latitudes)

	forecasts := make([]types.WeatherForecastProperties, 0, count)
	for i := 0; i < count; i++ {
		hasPrecipitationToday := false
		if i%2 == 0 {
			hasPrecipitationToday = true
		}
		forecasts = append(forecasts, types.WeatherForecastProperties{
			HasPrecipitationToday: hasPrecipitationToday,
			WeatherForecasts: []types.WeatherForecast{
				{
					Time:             "2024-09-06T00:00",
					Temperature2m:    0.1 + float64(i),
					Precipitation:    0.2 + float64(i),
					WindSpeed10m:     0.3 + float64(i),
					WindDirection10m: 0.4 + float64(i),
				},
				{
					Time:             "2024-09-06T01:00",
					Temperature2m:    0.1 + float64(i*10),
					Precipitation:    0.2 + float64(i*10),
					WindSpeed10m:     0.3 + float64(i*10),
					WindDirection10m: 0.4 + float64(i*10),
				},
			},
		})
	}

	return forecasts, nil
}
func (f *fakeWeatherAPI) GetElevation(ctx context.Context, latitude float64, longitude float64) (float64, error) {
	return latitude / longitude, nil
}

type fakeDB struct{}

func (f *fakeDB) CreatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error) {
	powerPlant.ID = 1
	powerPlant.CreatedAt = time.Now()
	return powerPlant, nil
}

func (f *fakeDB) UpdatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error) {
	if powerPlant.Version == 999 {
		return nil, sql.ErrNoRows
	}
	powerPlant.UpdatedAt = time.Now()
	powerPlant.Version = powerPlant.Version + 1
	return powerPlant, nil
}

func (f *fakeDB) GetPowerPlant(ctx context.Context, id int64) (*types.PowerPlant, error) {
	if id == 999 {
		return nil, sql.ErrNoRows
	}
	return &types.PowerPlant{
		ID:        id,
		Name:      "My Cool Power Plant",
		Latitude:  22.11,
		Longitude: 33.11,
		Version:   2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (f *fakeDB) GetPowerPlants(ctx context.Context, lastID int64, count int) ([]types.PowerPlant, error) {
	powerPlants := make([]types.PowerPlant, 0, count)
	for i := lastID + 1; i <= lastID+int64(count); i++ {
		powerPlants = append(powerPlants, types.PowerPlant{
			ID:        i,
			Name:      fmt.Sprintf("My Cool Power Plant %d", i),
			Latitude:  0.22 + float64(i*10),
			Longitude: 0.44 + float64(i*10),
			Version:   i,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return powerPlants, nil
}
