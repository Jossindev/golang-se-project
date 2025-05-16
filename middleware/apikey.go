package middleware

import (
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("api_key")
		expectedApiKey := os.Getenv("API_KEY")

		// For development, use a default if not set
		if expectedApiKey == "" {
			expectedApiKey = "special-key"
		}

		if apiKey == "" || apiKey != expectedApiKey {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(http.StatusUnauthorized, "Invalid or missing API key"))
			c.Abort()
			return
		}

		c.Next()
	}
}
