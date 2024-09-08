package usecase

import (
	"context"
	"database/sql"

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

// Usecase ...
type Usecase struct {
	weatherAPI weatherAPI
	db         db
}

// NewUsecase ...
func NewUsecase(weatherAPI weatherAPI, db db) *Usecase {
	return &Usecase{
		weatherAPI: weatherAPI,
		db:         db,
	}
}

// CreatePowerPlant validates and creates a new power plant.
func (u *Usecase) CreatePowerPlant(ctx context.Context, name string, lat float64, long float64) (*types.PowerPlant, error) {
	if name == "" {
		return nil, types.NewError("name is required").WithCode(types.ErrBadRequest)
	}
	if lat == 0 {
		return nil, types.NewError("latitude is required").WithCode(types.ErrBadRequest)
	}
	if long == 0 {
		return nil, types.NewError("longitude is required").WithCode(types.ErrBadRequest)
	}
	if lat > 90 || lat < -90 {
		return nil, types.NewError("latitude must be between -90 and 90").WithCode(types.ErrBadRequest)
	}
	if long > 180 || long < -180 {
		return nil, types.NewError("longitude must be between -180 and 180").WithCode(types.ErrBadRequest)
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
		return nil, types.NewError("id is required").WithCode(types.ErrBadRequest)
	}
	if lat != nil && (*lat > 90 || *lat < -90) {
		return nil, types.NewError("latitude must be between -90 and 90").WithCode(types.ErrBadRequest)
	}
	if long != nil && (*long > 180 || *long < -180) {
		return nil, types.NewError("longitude must be between -180 and 180").WithCode(types.ErrBadRequest)
	}

	powerPlant, err := u.db.GetPowerPlantForUpdate(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, types.NewError("id not found").WithCode(types.ErrBadRequest)
		default:
			// TODO: unhandled error should be logged as we're not passing the error to UI
			return nil, err
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
			return nil, types.NewError("id not found").WithCode(types.ErrBadRequest)
		default:
			// TODO: unhandled error should be logged as we're not passing the error to UI
			return nil, err
		}
	}
	return powerPlant, nil
}

// GetPowerPlant returns a power plant by ID.
func (u *Usecase) GetPowerPlant(ctx context.Context, id int64, forecastDays int) (*types.PowerPlant, error) {
	if id == 0 {
		return nil, types.NewError("id is required").WithCode(types.ErrBadRequest)
	}

	if _, ok := types.ValidForecastLengths[forecastDays]; !ok {
		return nil, types.NewError("invalid forecast days").WithCode(types.ErrBadRequest)
	}

	powerPlant, err := u.db.GetPowerPlant(ctx, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, types.NewError("id not found").WithCode(types.ErrBadRequest)
		default:
			// TODO: unhandled error should be logged as we're not passing the error to UI
			return nil, err
		}
	}

	forecast, err := u.weatherAPI.GetWeatherForecast(ctx, powerPlant.Latitude, powerPlant.Longitude, forecastDays)
	if err != nil {
		return nil, err
	}

	powerPlant.WeatherForecastProperties = *forecast

	elevations, err := u.weatherAPI.GetElevations(ctx, []float64{powerPlant.Latitude}, []float64{powerPlant.Longitude})
	if err != nil {
		return nil, err
	}

	powerPlant.Elevation = elevations[0]

	return powerPlant, nil
}

// GetPowerPlants returns a list of power plants.
// We use lastID to mark the last power plant ID we fetched instead of using offset to avoid performance issues when the table grows.
func (u *Usecase) GetPowerPlants(ctx context.Context, lastID int64, count int, forecastDays int) ([]types.PowerPlant, error) {
	if _, ok := types.ValidForecastLengths[forecastDays]; !ok {
		return nil, types.NewError("invalid forecast days").WithCode(types.ErrBadRequest)
	}

	powerPlants, err := u.db.GetPowerPlants(ctx, lastID, count)
	if err != nil {
		return nil, err
	}

	lats := make([]float64, 0, len(powerPlants))
	longs := make([]float64, 0, len(powerPlants))
	for _, powerPlant := range powerPlants {
		lats = append(lats, powerPlant.Latitude)
		longs = append(longs, powerPlant.Longitude)
	}

	forecasts, err := u.weatherAPI.GetWeatherForecasts(ctx, lats, longs, forecastDays)
	if err != nil {
		return nil, err
	}

	elevations, err := u.weatherAPI.GetElevations(ctx, lats, longs)
	if err != nil {
		return nil, err
	}

	for i, forecast := range forecasts {
		powerPlants[i].WeatherForecastProperties = forecast
		powerPlants[i].Elevation = elevations[i]
	}

	return powerPlants, nil
}
