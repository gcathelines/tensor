package main

import (
	"context"
	"testing"
	"time"

	"github.com/gcathelines/tensor-energy-case/internal/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/machinebox/graphql"
)

func TestE2EProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cl := graphql.NewClient("http://localhost:8080/query")
	createReq := createPowerPlantRequest("My Cool Power Plant", 1.1, 2.2)

	resp := map[string]types.PowerPlant{}
	err := cl.Run(ctx, createReq, &resp)
	if err != nil {
		t.Fatal(err)
	}

	getReq := getPowerPlantRequest(resp["createPowerPlant"].ID, 7)
	err = cl.Run(ctx, getReq, &resp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(types.PowerPlant{
		ID:        resp["createPowerPlant"].ID,
		Name:      "My Cool Power Plant",
		Latitude:  1.1,
		Longitude: 2.2,
	}, resp["powerPlant"], cmpopts.IgnoreFields(types.PowerPlant{}, "WeatherForecastProperties")); diff != "" {
		t.Fatalf("unexpected response (-want +got):\n%s", diff)
	}

	newName := "My Even Cooler Power Plant"
	updateReq := updatePowerPlantRequest(resp["createPowerPlant"].ID, &newName, nil, nil)
	err = cl.Run(ctx, updateReq, &resp)
	if err != nil {
		t.Fatal(err)
	}

	err = cl.Run(ctx, getReq, &resp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(types.PowerPlant{
		ID:        resp["createPowerPlant"].ID,
		Name:      newName,
		Latitude:  1.1,
		Longitude: 2.2,
	}, resp["powerPlant"], cmpopts.IgnoreFields(types.PowerPlant{}, "WeatherForecastProperties")); diff != "" {
		t.Fatalf("unexpected response (-want +got):\n%s", diff)
	}

	respPagination := map[string][]types.PowerPlant{}
	getPaginatedReq := getPaginatedPowerPlantsRequest(0, 2, 7)
	err = cl.Run(ctx, getPaginatedReq, &respPagination)
	if err != nil {
		t.Fatal(err)
	}

	if len(respPagination["powerPlants"]) != 2 {
		t.Fatalf("expected 2 power plants, got %d", len(respPagination["powerPlants"]))
	}
}

func getPowerPlantRequest(id int64, forecastDays int) *graphql.Request {
	query := `
    query GetPowerPlant($id: ID!, $forecastDays: Int = 7) {
      powerPlant(id: $id, forecastDays: $forecastDays) {
        id
        name
        latitude
        longitude
        elevation
        hasPrecipitationToday
        weatherForecasts(forecastDays: $forecastDays) {
          time
          temperature
          precipitation
          windSpeed
          windDirection
        }
      }
    }`

	req := graphql.NewRequest(query)
	req.Var("id", id)
	req.Var("forecastDays", forecastDays)

	return req
}

func getPaginatedPowerPlantsRequest(lastID int64, count, forecastDays int) *graphql.Request {
	query := `
    query GetPaginatedPowerPlants($lastID: Int64 = 0, $count: Int = 10, $forecastDays: Int = 7) {
      powerPlants(lastID: $lastID, count: $count, forecastDays: $forecastDays) {
        id
        name
        latitude
        longitude
        elevation
        hasPrecipitationToday
        weatherForecasts(forecastDays: $forecastDays) {
          time
          temperature
          precipitation
          windSpeed
          windDirection
        }
      }
    }`

	req := graphql.NewRequest(query)
	req.Var("lastID", lastID)
	req.Var("count", count)
	req.Var("forecastDays", forecastDays)

	return req
}

func createPowerPlantRequest(name string, latitude, longitude float64) *graphql.Request {
	mutation := `
    mutation CreatePowerPlant($input: CreatePowerPlantInput!) {
      createPowerPlant(input: $input) {
        id
        name
        latitude
        longitude
        elevation
        hasPrecipitationToday
      }
    }`

	req := graphql.NewRequest(mutation)
	req.Var("input", map[string]interface{}{
		"name":      name,
		"latitude":  latitude,
		"longitude": longitude,
	})

	return req
}

func updatePowerPlantRequest(id int64, name *string, latitude, longitude *float64) *graphql.Request {
	mutation := `
    mutation UpdatePowerPlant($input: UpdatePowerPlantInput!) {
      updatePowerPlant(input: $input) {
        id
        name
        latitude
        longitude
        elevation
        hasPrecipitationToday
      }
    }`

	input := map[string]interface{}{
		"id": id,
	}

	if name != nil {
		input["name"] = name
	}
	if latitude != nil {
		input["latitude"] = *latitude
	}
	if longitude != nil {
		input["longitude"] = *longitude
	}

	req := graphql.NewRequest(mutation)
	req.Var("input", input)

	return req
}
