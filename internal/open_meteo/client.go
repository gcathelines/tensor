package open_meteo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gcathelines/tensor-energy-case/config"
	"github.com/gcathelines/tensor-energy-case/internal/types"
)

// OpenMeteoClient is a client for the OpenMeteo API.
// Full documentation can be found at https://open-meteo.com/en/docs.
type OpenMeteoClient struct {
	apiURL     string
	httpClient *http.Client
}

// NewOpenMeteoClient creates a new OpenMeteoClient.
func NewOpenMeteoClient(cfg config.OpenMeteoConfig) *OpenMeteoClient {
	return &OpenMeteoClient{
		apiURL: cfg.APIURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// GetWeatherForecasts returns the weather forecast for a pair of latitude and longitude.
// Docs: https://open-meteo.com/en/docs/weather-api
func (c *OpenMeteoClient) GetWeatherForecast(ctx context.Context, latitude float64, longitude float64, forecastDays int) (*types.WeatherForecastProperties, error) {
	query := url.Values{
		"forecast_days": {fmt.Sprint(forecastDays)},
		"latitude":      {fmt.Sprint(latitude)},
		"longitude":     {fmt.Sprint(longitude)},
		"daily": {
			"precipitation_sum",
		},
		"hourly": {
			"temperature_2m",
			"precipitation",
			"wind_speed_10m",
			"wind_direction_10m",
		},
	}

	var forecast WeatherForecast
	err := c.doRequest("/v1/forecast", query, "GET", nil, &forecast)
	if err != nil {
		return nil, err
	}

	properties, err := forecast.ToProperties()
	if err != nil {
		return nil, err
	}

	return properties, nil
}

// GetWeatherForecasts returns the weather forecast for multiple pair latitude and longitude.
// Docs: https://open-meteo.com/en/docs/weather-api
func (c *OpenMeteoClient) GetWeatherForecasts(ctx context.Context, latitudes []float64, longitudes []float64, forecastDays int) ([]types.WeatherForecastProperties, error) {
	latsStr := make([]string, 0, len(latitudes))
	for _, lat := range latitudes {
		latsStr = append(latsStr, fmt.Sprint(lat))
	}

	longsStr := make([]string, 0, len(longitudes))
	for _, long := range longitudes {
		longsStr = append(longsStr, fmt.Sprint(long))
	}

	query := url.Values{
		"forecast_days": {fmt.Sprint(forecastDays)},
		"latitude":      latsStr,
		"longitude":     longsStr,
		"daily": {
			"precipitation_sum",
		},
		"hourly": {
			"temperature_2m",
			"precipitation",
			"wind_speed_10m",
			"wind_direction_10m",
		},
	}

	var forecasts []WeatherForecast
	err := c.doRequest("/v1/forecast", query, "GET", nil, &forecasts)
	if err != nil {
		return nil, err
	}

	properties := make([]types.WeatherForecastProperties, 0, len(forecasts))
	for _, forecast := range forecasts {
		props, err := forecast.ToProperties()
		if err != nil {
			return nil, err
		}

		properties = append(properties, *props)
	}

	return properties, nil
}

// GetElevation returns the elevation for the given latitude and longitude.
// Docs: https://open-meteo.com/en/docs/elevation-api
func (c *OpenMeteoClient) GetElevation(ctx context.Context, latitude float64, longitude float64) (float64, error) {
	query := url.Values{
		"latitude":  {fmt.Sprint(latitude)},
		"longitude": {fmt.Sprint(longitude)},
	}

	var elevation Elevation
	err := c.doRequest("/v1/elevation", query, "GET", nil, &elevation)
	if err != nil {
		return 0, err
	}

	if len(elevation.Elevation) != 1 {
		return 0, fmt.Errorf("unexpected number of elevations: %d", len(elevation.Elevation))
	}

	return elevation.Elevation[0], nil
}

// doRequest performs a request to the OpenMeteo API.
//  1. It constructs the URL with the given path and query parameters.
//  2. It creates a new HTTP request with the given method and body.
//  3. It sends the request and checks the response status code.
//     If the status code is not 200, it decodes the response body into an ErrorResponse object.
//  4. It decodes the response body into the given response object.
func (c *OpenMeteoClient) doRequest(
	path string,
	query url.Values,
	method string,
	body io.Reader,
	response any,
) error {
	reqURL, err := url.Parse(c.apiURL)
	if err != nil {
		return err
	}

	reqURL = reqURL.JoinPath(path)
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var (
			errorResponse ErrorResponse
			errorMsg      = fmt.Sprintf("unexpected status code: %d", res.StatusCode)
		)
		// We ignore the error here because we don't want to lose the original error.
		// Getting the error message would be nice, but it's not critical.
		_ = json.NewDecoder(res.Body).Decode(&errorResponse)
		if errorResponse.Error {
			errorMsg += fmt.Sprintf(", reason: %s", errorResponse.Reason)
		}

		return errors.New(errorMsg)
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	return nil
}
