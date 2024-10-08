package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/gcathelines/tensor-energy-case/internal/types"
)

// CreatePowerPlant is the resolver for the createPowerPlant field.
func (r *mutationResolver) CreatePowerPlant(ctx context.Context, input CreatePowerPlantInput) (*types.PowerPlant, error) {
	return r.usecase.CreatePowerPlant(ctx, input.Name, input.Latitude, input.Longitude)
}

// UpdatePowerPlant is the resolver for the updatePowerPlant field.
func (r *mutationResolver) UpdatePowerPlant(ctx context.Context, input UpdatePowerPlantInput) (*types.PowerPlant, error) {
	return r.usecase.UpdatePowerPlant(ctx, input.ID, input.Name, input.Latitude, input.Longitude)
}

// PowerPlant is the resolver for the powerPlant field.
func (r *queryResolver) PowerPlant(ctx context.Context, id int64, forecastDays *int) (*types.PowerPlant, error) {
	if forecastDays == nil {
		defaultForecastDays := 7
		forecastDays = &defaultForecastDays
	}

	return r.usecase.GetPowerPlant(ctx, id, *forecastDays)
}

// PowerPlants is the resolver for the powerPlants field.
func (r *queryResolver) PowerPlants(ctx context.Context, lastID *int64, count *int, forecastDays *int) ([]types.PowerPlant, error) {
	if lastID == nil {
		defaultLastID := int64(0)
		lastID = &defaultLastID
	}

	if count == nil {
		defaultCount := 10
		count = &defaultCount
	}

	if forecastDays == nil {
		defaultForecastDays := 7
		forecastDays = &defaultForecastDays
	}

	return r.usecase.GetPowerPlants(ctx, *lastID, *count, *forecastDays)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
