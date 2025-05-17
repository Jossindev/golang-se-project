package handlers

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"awesomeProject/services"
	"awesomeProject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Subscribe(c *gin.Context) {
	var request models.SubscriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingSubscription models.Subscription
	result := db.DB.Where("email = ? AND city = ?", request.Email, request.City).First(&existingSubscription)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already subscribed to this city's weather updates"})
		return
	}

	token, err := utils.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	subscription := models.Subscription{
		Email:     request.Email,
		City:      request.City,
		Frequency: request.Frequency,
		Confirmed: false,
		Token:     token,
	}

	if err := db.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	emailService := services.NewEmailService()
	if err := emailService.SendConfirmationEmail(subscription.Email, subscription.City, subscription.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription created. Please check your email to confirm."})
}

func ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	var subscription models.Subscription
	result := db.DB.Where("token = ?", token).First(&subscription)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	subscription.Confirmed = true

	newToken, err := utils.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	subscription.Token = newToken

	if err := db.DB.Save(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm subscription"})
		return
	}

	emailService := services.NewEmailService()
	if err := emailService.SendUnsubscribeEmail(subscription.Email, subscription.City, subscription.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send unsubscribe instructions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}

func Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	var subscription models.Subscription
	result := db.DB.Where("token = ?", token).First(&subscription)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	if err := db.DB.Delete(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}
