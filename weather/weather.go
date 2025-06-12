package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://api.weather.gov"

	// Temperature thresholds in Fahrenheit
	HotThreshold      = 80.0
	ColdThreshold     = 50.0
	ModerateThreshold = 51.0 // Lower bound for moderate
)

// TemperatureCategory represents the temperature characterization
type TemperatureCategory string

const (
	Hot      TemperatureCategory = "hot"
	Cold     TemperatureCategory = "cold"
	Moderate TemperatureCategory = "moderate"
)

// WeatherResponse represents the response from our API
type WeatherResponse struct {
	Forecast     string              `json:"forecast"`
	Temperature  float64             `json:"temperature"`
	Category     TemperatureCategory `json:"category"`
	ErrorMessage string              `json:"error,omitempty"`
}

// NWSResponse represents the National Weather Service API response
type NWSResponse struct {
	Properties struct {
		Periods []struct {
			Temperature   float64 `json:"temperature"`
			ShortForecast string  `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

// gridResponse represents the response from the grid endpoint
type gridResponse struct {
	Properties struct {
		ForecastURL string `json:"forecast"`
	} `json:"properties"`
}

// buildGridURL constructs the grid endpoint URL safely
func buildGridURL(lat, lon float64) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Format coordinates to 4 decimal places for consistency
	latStr := fmt.Sprintf("%.4f", lat)
	lonStr := fmt.Sprintf("%.4f", lon)

	// Build the path
	base.Path = fmt.Sprintf("/points/%s,%s", latStr, lonStr)

	return base.String(), nil
}

// fetchGridData gets the grid data for the given coordinates
func fetchGridData(client *http.Client, lat, lon float64) (*gridResponse, error) {
	gridURL, err := buildGridURL(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("error building grid URL: %w", err)
	}

	req, err := http.NewRequest("GET", gridURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "WeatherService/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var gridResp gridResponse
	if err := json.Unmarshal(body, &gridResp); err != nil {
		return nil, fmt.Errorf("error parsing grid response: %w", err)
	}

	return &gridResp, nil
}

// fetchForecast gets the forecast data from the given URL
func fetchForecast(client *http.Client, forecastURL string) (*NWSResponse, error) {
	// Validate the forecast URL
	url, err := url.Parse(forecastURL)
	if err != nil {
		return nil, fmt.Errorf("invalid forecast URL: %w", err)
	}

	// Ensure the forecast URL is from the same domain
	if url.Host != "api.weather.gov" {
		return nil, fmt.Errorf("invalid forecast URL domain")
	}

	resp, err := client.Get(forecastURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected forecast status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading forecast response: %w", err)
	}

	var nwsResp NWSResponse
	if err := json.Unmarshal(body, &nwsResp); err != nil {
		return nil, fmt.Errorf("error parsing forecast response: %w", err)
	}

	if len(nwsResp.Properties.Periods) == 0 {
		return nil, fmt.Errorf("no forecast periods available")
	}

	return &nwsResp, nil
}

// GetWeather fetches weather data for the given coordinates
func GetWeather(lat, lon float64) (*WeatherResponse, error) {
	// Validate coordinates
	if !IsValidCoordinates(lat, lon) {
		return nil, fmt.Errorf("invalid coordinates: latitude must be between -90 and 90, longitude between -180 and 180")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Get grid data
	gridResp, err := fetchGridData(client, lat, lon)
	if err != nil {
		return nil, err
	}

	// Get forecast data
	forecastResp, err := fetchForecast(client, gridResp.Properties.ForecastURL)
	if err != nil {
		return nil, err
	}

	// Get today's forecast (first period)
	today := forecastResp.Properties.Periods[0]

	return &WeatherResponse{
		Forecast:    today.ShortForecast,
		Temperature: today.Temperature,
		Category:    categorizeTemperature(today.Temperature),
	}, nil
}

// IsValidCoordinates checks if the given coordinates are valid
func IsValidCoordinates(lat, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

// categorizeTemperature determines the temperature category
func categorizeTemperature(temp float64) TemperatureCategory {
	switch {
	case temp >= HotThreshold:
		return Hot
	case temp <= ColdThreshold:
		return Cold
	default:
		return Moderate
	}
}
