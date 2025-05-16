package middleware

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(http.StatusUnauthorized, "Authorization header is required"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(http.StatusUnauthorized, "Invalid token: "+err.Error()))
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

func UserOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathUsername := c.Param("username")
		loggedInUsername, exists := c.Get("username")

		if !exists || pathUsername != loggedInUsername.(string) {
			c.JSON(http.StatusForbidden, models.ErrorResponse(http.StatusForbidden, "You can only access your own user resources"))
			c.Abort()
			return
		}

		c.Next()
	}
}
