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
					Time:          "2024-09-06T00:00",
					Temperature:   0.1 + float64(i),
					Precipitation: 0.2 + float64(i),
					WindSpeed:     0.3 + float64(i),
					WindDirection: 0.4 + float64(i),
				},
				{
					Time:          "2024-09-06T01:00",
					Temperature:   0.1 + float64(i*10),
					Precipitation: 0.2 + float64(i*10),
					WindSpeed:     0.3 + float64(i*10),
					WindDirection: 0.4 + float64(i*10),
				},
			},
		})
	}

	return forecasts, nil
}
func (f *fakeWeatherAPI) GetElevations(ctx context.Context, latitude []float64, longitude []float64) ([]float64, error) {
	res := make([]float64, 0, len(latitude))
	for i := 0; i < len(latitude); i++ {
		res = append(res, latitude[i]/longitude[i])
	}

	return res, nil
}

type fakeDB struct{}

func (f *fakeDB) CreatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error) {
	powerPlant.ID = 1
	powerPlant.CreatedAt = time.Now()
	return powerPlant, nil
}

func (f *fakeDB) UpdatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error) {
	if powerPlant.ID == 999 {
		return nil, sql.ErrNoRows
	}
	powerPlant.UpdatedAt = time.Now()
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (f *fakeDB) GetPowerPlantForUpdate(ctx context.Context, id int64) (*types.PowerPlant, error) {
	if id == 999 {
		return nil, sql.ErrNoRows
	}
	return &types.PowerPlant{
		ID:        id,
		Name:      "My Cool Power Plant",
		Latitude:  22.11,
		Longitude: 33.11,
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return powerPlants, nil
}
