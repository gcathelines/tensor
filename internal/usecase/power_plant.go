package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/gcathelines/tensor-energy-case/internal/database"
	"github.com/gcathelines/tensor-energy-case/internal/open_meteo"
	"github.com/gcathelines/tensor-energy-case/internal/types"
)

var _ weatherAPI = (*open_meteo.OpenMeteoClient)(nil)

type weatherAPI interface {
	GetWeatherForecast(ctx context.Context, latitudes float64, longitudes float64, forecastDays int) (*types.WeatherForecastProperties, error)
	GetWeatherForecasts(ctx context.Context, latitudes []float64, longitudes []float64, forecastDays int) ([]types.WeatherForecastProperties, error)
	GetElevations(ctx context.Context, latitude []float64, longitude []float64) ([]float64, error)
}

var _ db = (*database.Database)(nil)

type db interface {
	CreatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error)
	UpdatePowerPlant(ctx context.Context, powerPlant *types.PowerPlant) (*types.PowerPlant, error)
	GetPowerPlant(ctx context.Context, id int64) (*types.PowerPlant, error)
	GetPowerPlantForUpdate(ctx context.Context, id int64) (*types.PowerPlant, error)
	GetPowerPlants(ctx context.Context, lastID int64, count int) ([]types.PowerPlant, error)
}

// Usecase represents the usecase of the service.
type Usecase struct {
	weatherAPI weatherAPI
	db         db
	logger     *log.Logger
}

// NewUsecase creates a new usecase.
func NewUsecase(weatherAPI weatherAPI, db db) *Usecase {

	return &Usecase{
		weatherAPI: weatherAPI,
		db:         db,
		logger:     log.Default(),
	}
}

// CreatePowerPlant validates and creates a new power plant.
func (u *Usecase) CreatePowerPlant(ctx context.Context, name string, lat float64, long float64) (*types.PowerPlant, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if lat == 0 {
		return nil, errors.New("latitude is required")
	}
	if long == 0 {
		return nil, errors.New("longitude is required")
	}
	if lat > 90 || lat < -90 {
		return nil, types.ErrInvalidLatitude
	}
	if long > 180 || long < -180 {
		return nil, types.ErrInvalidLongitude
	}

	return u.db.CreatePowerPlant(ctx, &types.PowerPlant{
		Name:      name,
		Latitude:  lat,
		Longitude: long,
	})
}

// UpdatePowerPlant updates a power plant by ID.
// We will use pessimistic lock to avoid write conflicts.
func (u *Usecase) UpdatePowerPlant(ctx context.Context, id int64, name *string, lat *float64, long *float64) (*types.PowerPlant, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	if lat != nil && (*lat > 90 || *lat < -90) {
		return nil, types.ErrInvalidLatitude
	}
	if long != nil && (*long > 180 || *long < -180) {
		return nil, types.ErrInvalidLongitude
	}

	powerPlant, err := u.db.GetPowerPlantForUpdate(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("id not found")
		default:
			u.logger.Printf("error getting power plant for update: %v", err)
			return nil, types.ErrInternal
		}
	}

	if lat != nil {
		powerPlant.Latitude = *lat
	}
	if long != nil {
		powerPlant.Longitude = *long
	}
	if name != nil {
		powerPlant.Name = *name
	}

	powerPlant, err = u.db.UpdatePowerPlant(ctx, powerPlant)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("id not found")
		default:
			u.logger.Printf("error updating power plant: %v", err)
			return nil, types.ErrInternal
		}
	}
	return powerPlant, nil
}

// GetPowerPlant returns a power plant by ID.
func (u *Usecase) GetPowerPlant(ctx context.Context, id int64, forecastDays int) (*types.PowerPlant, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}

	if _, ok := types.ValidForecastLengths[forecastDays]; !ok {
		return nil, types.ErrInvalidForecastDay
	}

	powerPlant, err := u.db.GetPowerPlant(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("id not found")
		default:
			u.logger.Printf("error getting power plant: %v", err)
			return nil, types.ErrInternal
		}
	}

	forecast, err := u.weatherAPI.GetWeatherForecast(ctx, powerPlant.Latitude, powerPlant.Longitude, forecastDays)
	if err != nil {
		u.logger.Printf("error getting weather forecast: %v", err)
		return nil, types.ErrInternal
	}

	powerPlant.WeatherForecastProperties = *forecast

	elevations, err := u.weatherAPI.GetElevations(ctx, []float64{powerPlant.Latitude}, []float64{powerPlant.Longitude})
	if err != nil {
		u.logger.Printf("error getting elevation: %v", err)
		return nil, types.ErrInternal
	}

	powerPlant.Elevation = elevations[0]

	return powerPlant, nil
}

// GetPowerPlants returns a list of power plants.
// We use lastID to mark the last power plant ID we fetched instead of using offset to avoid performance issues when the table grows.
func (u *Usecase) GetPowerPlants(ctx context.Context, lastID int64, count int, forecastDays int) ([]types.PowerPlant, error) {
	if _, ok := types.ValidForecastLengths[forecastDays]; !ok {
		return nil, types.ErrInvalidForecastDay
	}

	powerPlants, err := u.db.GetPowerPlants(ctx, lastID, count)
	if err != nil {
		return nil, err
	}

	if len(powerPlants) == 0 {
		return powerPlants, nil
	}

	lats := make([]float64, 0, len(powerPlants))
	longs := make([]float64, 0, len(powerPlants))
	for _, powerPlant := range powerPlants {
		lats = append(lats, powerPlant.Latitude)
		longs = append(longs, powerPlant.Longitude)
	}

	forecasts, err := u.weatherAPI.GetWeatherForecasts(ctx, lats, longs, forecastDays)
	if err != nil {
		u.logger.Printf("error getting weather forecasts: %v", err)
		return nil, types.ErrInternal
	}

	elevations, err := u.weatherAPI.GetElevations(ctx, lats, longs)
	if err != nil {
		u.logger.Printf("error getting elevations: %v", err)
		return nil, types.ErrInternal
	}

	for i, forecast := range forecasts {
		powerPlants[i].WeatherForecastProperties = forecast
		powerPlants[i].Elevation = elevations[i]
	}

	return powerPlants, nil
}
