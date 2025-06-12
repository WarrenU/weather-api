package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather/weather"

	"github.com/gorilla/mux"
)

func TestWeatherHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		checkResponse  func(*testing.T, *weather.WeatherResponse)
	}{
		{
			name:           "Missing parameters",
			query:          "",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp *weather.WeatherResponse) {
				if resp.ErrorMessage != "latitude and longitude are required" {
					t.Errorf("Expected error message about missing parameters, got: %v", resp.ErrorMessage)
				}
			},
		},
		{
			name:           "Invalid latitude",
			query:          "lat=invalid&lon=-74.0060",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp *weather.WeatherResponse) {
				if resp.ErrorMessage != "invalid latitude format" {
					t.Errorf("Expected error message about invalid latitude, got: %v", resp.ErrorMessage)
				}
			},
		},
		{
			name:           "Invalid longitude",
			query:          "lat=40.7128&lon=invalid",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp *weather.WeatherResponse) {
				if resp.ErrorMessage != "invalid longitude format" {
					t.Errorf("Expected error message about invalid longitude, got: %v", resp.ErrorMessage)
				}
			},
		},
		{
			name:           "Out of range coordinates",
			query:          "lat=91.0&lon=-74.0060",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp *weather.WeatherResponse) {
				if resp.ErrorMessage != "coordinates out of valid range" {
					t.Errorf("Expected error message about invalid coordinates, got: %v", resp.ErrorMessage)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			req := httptest.NewRequest("GET", "/weather?"+tt.query, nil)
			w := httptest.NewRecorder()

			// Create a new router for testing
			r := mux.NewRouter()
			r.Handle("/weather", coordinateValidator(http.HandlerFunc(weatherHandler))).Methods("GET")

			// Serve the request
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			var response weather.WeatherResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			tt.checkResponse(t, &response)
		})
	}
}

// TestWeatherHandlerWithValidCoordinates tests the handler with valid coordinates
func TestWeatherHandlerWithValidCoordinates(t *testing.T) {
	// Create a request with valid coordinates
	req := httptest.NewRequest("GET", "/weather?lat=40.7128&lon=-74.0060", nil)
	w := httptest.NewRecorder()

	// Create a new router for testing
	r := mux.NewRouter()
	r.Handle("/weather", coordinateValidator(http.HandlerFunc(weatherHandler))).Methods("GET")

	// Serve the request
	r.ServeHTTP(w, req)

	// Check if the request was processed
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 200 or 500, got %d", w.Code)
	}

	var response weather.WeatherResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// If we got a successful response, verify the structure
	if w.Code == http.StatusOK {
		if response.Forecast == "" {
			t.Error("Expected non-empty forecast")
		}
		if response.Category == "" {
			t.Error("Expected non-empty temperature category")
		}
	}
}
