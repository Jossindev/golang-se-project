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

	// User routes
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

	// Pet routes
	// Routes protected by OAuth (JWT in our case)
	petAuth := router.Group("")
	petAuth.Use(middleware.AuthMiddleware())
	{
		petAuth.POST("/pet", handlers.AddPet)
		petAuth.PUT("/pet", handlers.UpdatePet)
		petAuth.GET("/pet/findByStatus", handlers.FindPetsByStatus)
		petAuth.GET("/pet/findByTags", handlers.FindPetsByTags)
		petAuth.POST("/pet/:petId", handlers.UpdatePetWithForm)
		petAuth.DELETE("/pet/:petId", handlers.DeletePet)
		petAuth.POST("/pet/:petId/uploadImage", handlers.UploadFile)
	}

	// Routes that can use either API key or OAuth
	router.GET("/pet/:petId", middleware.ApiKeyAuth(), handlers.GetPetById)

	router.Run(":8080")
}
