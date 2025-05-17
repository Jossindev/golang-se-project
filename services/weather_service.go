package services

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type WeatherService struct {
	APIKey  string
	BaseURL string
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		APIKey:  utils.GetEnv("API_KEY", ""), //Set API_KY with token which will be provided by Weather API
		BaseURL: "http://api.weatherapi.com/v1/current.json",
	}
}

func (s *WeatherService) GetWeather(city string) (*models.Weather, error) {
	reqURL := fmt.Sprintf("%s?key=%s&q=%s", s.BaseURL, s.APIKey, url.QueryEscape(city))

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("city not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API error: %s", resp.Status)
	}

	var apiResponse struct {
		Current struct {
			TempC     float64 `json:"temp_c"`
			Humidity  float64 `json:"humidity"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	weather := &models.Weather{
		Temperature: apiResponse.Current.TempC,
		Humidity:    apiResponse.Current.Humidity,
		Description: apiResponse.Current.Condition.Text,
	}

	return weather, nil
}
