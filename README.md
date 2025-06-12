# Weather Service

A Go-based HTTP server that provides weather forecasts using the National Weather Service API.

## Features

- Accepts lat and lon coordinates as query parameters
- Returns short forecast for the specified location
- Provides temperature characterization (hot, cold, or moderate)
- Input validation and error handling
- Graceful error handling and resource cleanup

## Prerequisites

- Go 1.21

## Installation

1. Clone the repository:
```bash
git clone https://github.com/WarrenU/weather-api
cd weather
```

2. Install dependencies:
```bash
go mod download
```

## Building

Build the application:
```bash
go build
```

## Testing

Run the unit tests:
```bash
go test ./...
```

This will run all tests in the project, including:
- Coordinate validation tests
- Temperature categorization tests
- HTTP handler tests

## Running

Start the server:
```bash
./weather
```

The server will start on port 8080.

## API Usage

### Get Weather Forecast

**Endpoint:** `GET /weather`

**Query Parameters:**
- `lat`: Latitude (required, -90 to 90)
- `lon`: Longitude (required, -180 to 180)

**Example Request:**
```bash
curl "http://localhost:8080/weather?lat=40.7128&lon=-74.0060"
```

**Example Response:**
```json
{
    "forecast": "Partly Cloudy",
    "temperature": 75.0,
    "category": "moderate"
}
```

### Finding Coordinates

You can find latitude and longitude coordinates for any location using these methods:

1. **Google Maps**
   - Go to [Google Maps](https://www.google.com/maps)
   - Search for your location
   - Right-click on the map and coordinates are the first option that shows up. If you click them they are copy-able.

2. **LatLong.net**
   - Visit [LatLong.net](https://www.latlong.net/)
   - Search for your city or location
   - Get the coordinates directly

3. **Example Coordinates:**
   - New York City: 40.7128, -74.0060
   - Los Angeles: 34.0522, -118.2437
   - London: 51.5074, -0.1278
   - Tokyo: 35.6762, 139.6503

Note: Coordinates should be provided in decimal degrees format, where:
- Latitude ranges from -90 (South) to 90 (North)
- Longitude ranges from -180 (West) to 180 (East)

## Error Handling

The API returns appropriate HTTP status codes and error messages:
- 400 Bad Request: Invalid or missing parameters
- 500 Internal Server Error: Server-side errors

## Temperature Categories in Fahrenheit

- Hot: >= 80
- Cold: <= 50
- Moderate: 51 to 79