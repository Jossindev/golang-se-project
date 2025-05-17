package services

import (
	"awesomeProject/models"
	"awesomeProject/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherService_GetWeather(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		city := query.Get("q")

		if query.Get("key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid API key"}`))
			return
		}

		switch city {
		case "London":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"current": {
					"temp_c": 15.5,
					"humidity": 70,
					"condition": {
						"text": "Partly cloudy"
					}
				}
			}`))
		case "NonExistentCity":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "City not found"}`))
		case "ErrorCity":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal server error"}`))
		case "InvalidJSON":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{invalid json}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"current": {
					"temp_c": 20.0,
					"humidity": 60,
					"condition": {
						"text": "Sunny"
					}
				}
			}`))
		}
	}))
	defer server.Close()

	service := &services.WeatherService{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
	}

	testCases := []struct {
		name          string
		city          string
		expectedError bool
		expected      *models.Weather
	}{
		{
			name:          "Valid city returns weather data",
			city:          "London",
			expectedError: false,
			expected: &models.Weather{
				Temperature: 15.5,
				Humidity:    70,
				Description: "Partly cloudy",
			},
		},
		{
			name:          "City not found returns error",
			city:          "NonExistentCity",
			expectedError: true,
			expected:      nil,
		},
		{
			name:          "API error returns error",
			city:          "ErrorCity",
			expectedError: true,
			expected:      nil,
		},
		{
			name:          "Invalid JSON returns error",
			city:          "InvalidJSON",
			expectedError: true,
			expected:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			weather, err := service.GetWeather(tc.city)

			if tc.expectedError && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Did not expect an error but got: %v", err)
			}

			if !tc.expectedError {
				if weather == nil {
					t.Fatal("Expected weather data but got nil")
				}

				if weather.Temperature != tc.expected.Temperature {
					t.Errorf("Expected temperature %v, got %v", tc.expected.Temperature, weather.Temperature)
				}

				if weather.Humidity != tc.expected.Humidity {
					t.Errorf("Expected humidity %v, got %v", tc.expected.Humidity, weather.Humidity)
				}

				if weather.Description != tc.expected.Description {
					t.Errorf("Expected description %q, got %q", tc.expected.Description, weather.Description)
				}
			}
		})
	}
}
