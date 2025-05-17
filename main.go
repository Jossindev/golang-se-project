package main

import (
	"awesomeProject/db"
	"awesomeProject/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {

	db.InitDB()

	router := gin.Default()

	currentDir, _ := os.Getwd()
	indexPath := currentDir + "/public/index.html"

	router.LoadHTMLFiles(indexPath)
	router.Static("/static", "./public")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

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
