package main

import (
	"awesomeProject/db"
	"awesomeProject/handlers"
	"awesomeProject/middleware"
	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	router := gin.Default()

	// Public routes
	router.POST("/user", handlers.CreateUser)
	router.POST("/user/createWithList", handlers.CreateUsersWithList)
	router.GET("/user/login", handlers.LoginUser)
	router.GET("/user/logout", handlers.LogoutUser)

	// Protected routes
	authorized := router.Group("")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Basic authorized routes
		authorized.GET("/user/:username", handlers.GetUserByName)

		// Routes that require ownership verification
		ownership := authorized.Group("")
		ownership.Use(middleware.UserOwnershipMiddleware())
		{
			ownership.PUT("/user/:username", handlers.UpdateUser)
			ownership.DELETE("/user/:username", handlers.DeleteUser)
		}
	}

	router.Run(":8080")
}
