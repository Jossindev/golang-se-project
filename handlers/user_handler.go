package handlers

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func CreateUser1(c *gin.Context) {
	// 2. Parse request body based on content type
	var user models.User
	contentType := c.ContentType()

	switch contentType {
	case "application/json":
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid JSON: "+err.Error()))
			return
		}
	case "application/xml":
		if err := c.ShouldBindXML(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid XML: "+err.Error()))
			return
		}
	case "application/x-www-form-urlencoded":
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid form data: "+err.Error()))
			return
		}
	default:
		c.JSON(http.StatusUnsupportedMediaType, models.ErrorResponse(http.StatusUnsupportedMediaType,
			"Content-Type must be application/json, application/xml, or application/x-www-form-urlencoded"))
		return
	}

	// 3. Validate required fields
	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Username and password are required"))
		return
	}

	// 4. Check if username already exists
	var existingUser models.User
	if result := db.DB.Where("username = ?", user.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Username already exists"))
		return
	}

	// 5. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to hash password"))
		return
	}
	user.Password = string(hashedPassword)

	// 6. Create user
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to create user: "+err.Error()))
		return
	}

	user.Password = ""

	// 8. Return appropriate content type
	switch contentType {
	case "application/xml":
		c.XML(http.StatusOK, user)
	default:
		c.JSON(http.StatusOK, user)

	}
}
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to hash password"))
		return
	}
	user.Password = string(hashedPassword)

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to create user"))
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func CreateUsersWithList(c *gin.Context) {
	var users []models.User
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	for i := range users {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
		if err != nil {
			continue
		}
		users[i].Password = string(hashedPassword)

		db.DB.Create(&users[i])
	}

	if len(users) > 0 {
		lastUser := users[len(users)-1]
		c.JSON(http.StatusOK, lastUser.ToResponse())
	} else {
		c.JSON(http.StatusOK, models.UserResponse{})
	}
}

func LoginUser(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Username and password are required"))
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid username/password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid username/password"))
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to generate token"))
		return
	}

	c.Header("X-Rate-Limit", "1000")
	c.Header("X-Expires-After", time.Now().Add(24*time.Hour).Format(time.RFC3339))

	c.JSON(http.StatusOK, token)
}

func LogoutUser(c *gin.Context) {
	// In a stateless JWT-based auth system, logout is handled client-side
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func GetUserByName(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "User not found"))
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func UpdateUser(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "User not found"))
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if updatedUser.Password != "" {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to hash password"))
			return
		}
		user.Password = string(hashedPassword)
	}

	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	user.Phone = updatedUser.Phone
	user.UserStatus = updatedUser.UserStatus

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to update user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "User not found"))
		return
	}

	if err := db.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to delete user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
