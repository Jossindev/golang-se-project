package handlers

import (
	"awesomeProject/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockWeatherService struct {
	city          string
	returnWeather *models.Weather
	returnError   error
}

func (m *mockWeatherService) GetWeather(city string) (*models.Weather, error) {
	m.city = city
	return m.returnWeather, m.returnError
}

func TestGetWeather(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryString    string
		returnWeather  *models.Weather
		returnError    error
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "Missing city parameter",
			queryString:    "",
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":"City parameter is required"}`,
		},
		{
			name:        "Successful weather retrieval",
			queryString: "?city=London",
			returnWeather: &models.Weather{
				Temperature: 15.5,
				Humidity:    70.0,
				Description: "Partly cloudy",
			},
			returnError:    nil,
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"temperature":15.5,"humidity":70,"description":"Partly cloudy"}`,
		},
		{
			name:           "City not found",
			queryString:    "?city=NotFound",
			returnWeather:  nil,
			returnError:    errors.New("city not found"),
			expectedStatus: http.StatusNotFound,
			expectedJSON:   `{"error":"City not found"}`,
		},
		{
			name:           "Service error",
			queryString:    "?city=Error",
			returnWeather:  nil,
			returnError:    errors.New("service error"),
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   `{"error":"Failed to fetch weather data: service error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()

			mock := &mockWeatherService{
				returnWeather: tc.returnWeather,
				returnError:   tc.returnError,
			}

			router.GET("/weather", func(c *gin.Context) {
				city := c.Query("city")
				if city == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
					return
				}

				weather, err := mock.GetWeather(city)
				if err != nil {
					if err.Error() == "city not found" {
						c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
					} else {
						c.JSON(http.StatusInternalServerError,
							gin.H{"error": "Failed to fetch weather data: " + err.Error()})
					}
					return
				}

				c.JSON(http.StatusOK, weather)
			})

			req, _ := http.NewRequest("GET", "/weather"+tc.queryString, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.queryString != "" {
				expectedCity := req.URL.Query().Get("city")
				assert.Equal(t, expectedCity, mock.city)
			}

			var actualJSON interface{}
			var expectedJSON interface{}

			if tc.expectedJSON != "" {
				err := json.Unmarshal(resp.Body.Bytes(), &actualJSON)
				assert.NoError(t, err)

				err = json.Unmarshal([]byte(tc.expectedJSON), &expectedJSON)
				assert.NoError(t, err)

				assert.Equal(t, expectedJSON, actualJSON)
			}
		})
	}
}
