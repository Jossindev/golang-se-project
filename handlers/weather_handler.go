package handlers

import (
	"awesomeProject/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	weatherService := services.NewWeatherService()

	weather, err := weatherService.GetWeather(city)
	if err != nil {
		if err.Error() == "city not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather data: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, weather)
}
