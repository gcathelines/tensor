// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

type CreatePowerPlantInput struct {
	// Name of the power plant
	Name string `json:"name"`
	// Latitude in degrees
	Latitude float64 `json:"latitude"`
	// Longitude in degrees
	Longitude float64 `json:"longitude"`
}

type Mutation struct {
}

type Query struct {
}

type UpdatePowerPlantInput struct {
	// ID of the power plant
	ID int64 `json:"id"`
	// Latest version of the power plant
	Version int64 `json:"version"`
	// Name of the power plant
	Name string `json:"name"`
	// Latitude in degrees
	Latitude float64 `json:"latitude"`
	// Longitude in degrees
	Longitude float64 `json:"longitude"`
}
