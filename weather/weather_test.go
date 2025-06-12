package weather

import (
	"testing"
)

func TestIsValidCoordinates(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lon      float64
		expected bool
	}{
		{"Valid coordinates", 40.7128, -74.0060, true},
		{"Valid coordinates at boundaries", 90.0, 180.0, true},
		{"Valid coordinates at negative boundaries", -90.0, -180.0, true},
		{"Invalid latitude too high", 91.0, -74.0060, false},
		{"Invalid latitude too low", -91.0, -74.0060, false},
		{"Invalid longitude too high", 40.7128, 181.0, false},
		{"Invalid longitude too low", 40.7128, -181.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCoordinates(tt.lat, tt.lon)
			if result != tt.expected {
				t.Errorf("isValidCoordinates(%v, %v) = %v; want %v",
					tt.lat, tt.lon, result, tt.expected)
			}
		})
	}
}

func TestCategorizeTemperature(t *testing.T) {
	tests := []struct {
		name     string
		temp     float64
		expected TemperatureCategory
	}{
		// Hot temperatures
		{"Hot temperature above threshold", HotThreshold + 5.0, Hot},
		{"Hot temperature at threshold", HotThreshold, Hot},
		{"Hot temperature just above threshold", HotThreshold + 0.1, Hot},

		// Cold temperatures
		{"Cold temperature below threshold", ColdThreshold - 5.0, Cold},
		{"Cold temperature at threshold", ColdThreshold, Cold},
		{"Cold temperature just below threshold", ColdThreshold - 0.1, Cold},

		// Moderate temperatures
		{"Moderate temperature middle", (HotThreshold + ColdThreshold) / 2, Moderate},
		{"Moderate temperature at lower bound", ModerateThreshold, Moderate},
		{"Moderate temperature just above lower bound", ModerateThreshold + 0.1, Moderate},
		{"Moderate temperature just below hot threshold", HotThreshold - 0.1, Moderate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := categorizeTemperature(tt.temp)
			if result != tt.expected {
				t.Errorf("categorizeTemperature(%v) = %v; want %v",
					tt.temp, result, tt.expected)
			}
		})
	}
}
