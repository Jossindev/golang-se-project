package main

import (
	"awesomeProject/db"
	"awesomeProject/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db.InitDB()

	// Set up Gin router
	r := gin.Default()

	// User routes - all routes are now public without auth
	r.POST("/user", handlers.CreateUser)
	r.POST("/user/createWithList", handlers.CreateUsersWithList)
	r.GET("/user/login", handlers.LoginUser)
	r.GET("/user/logout", handlers.LogoutUser)
	r.GET("/user/:username", handlers.GetUserByName)
	r.PUT("/user/:username", handlers.UpdateUser)
	r.DELETE("/user/:username", handlers.DeleteUser)

	r.Run(":8080") // Listen on port 8080
}
