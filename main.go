package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"weather/weather"

	"github.com/gorilla/mux"
)

// coordinateValidator middleware validates coordinate parameters
func coordinateValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Get query parameters
		latStr := r.URL.Query().Get("lat")
		lonStr := r.URL.Query().Get("lon")

		// Validate required parameters
		if latStr == "" || lonStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(weather.WeatherResponse{
				ErrorMessage: "latitude and longitude are required",
			})
			return
		}

		// Parse coordinates
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(weather.WeatherResponse{
				ErrorMessage: "invalid latitude format",
			})
			return
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(weather.WeatherResponse{
				ErrorMessage: "invalid longitude format",
			})
			return
		}

		// Validate coordinate ranges
		if !weather.IsValidCoordinates(lat, lon) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(weather.WeatherResponse{
				ErrorMessage: "coordinates out of valid range",
			})
			return
		}

		// Store validated coordinates in request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "lat", lat)
		ctx = context.WithValue(ctx, "lon", lon)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// Get coordinates from context
	lat := r.Context().Value("lat").(float64)
	lon := r.Context().Value("lon").(float64)

	// Get weather data
	weatherData, err := weather.GetWeather(lat, lon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(weather.WeatherResponse{
			ErrorMessage: err.Error(),
		})
		return
	}

	// Return successful response
	json.NewEncoder(w).Encode(weatherData)
}

func main() {
	// Create router
	r := mux.NewRouter()

	// Define routes with middleware
	r.Handle("/weather", coordinateValidator(http.HandlerFunc(weatherHandler))).Methods("GET")

	// Add logging middleware
	r.Use(loggingMiddleware)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs all requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
