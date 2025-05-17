package main

import (
	"awesomeProject/db"
	"awesomeProject/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	router := gin.Default()

	api := router.Group("/api")
	{
		// Weather routes
		api.GET("/weather", handlers.GetWeather)

		// Subscription routes
		api.POST("/subscribe", handlers.Subscribe)
		api.GET("/confirm/:token", handlers.ConfirmSubscription)
		api.GET("/unsubscribe/:token", handlers.Unsubscribe)
	}

	router.Run(":8080")
}
