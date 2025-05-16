package handlers

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"fmt"
	"github.com/gin-gonic/gin"

	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func AddPet(c *gin.Context) {
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if userID, exists := c.Get("userID"); exists {
		pet.OwnerID = userID.(uint)
	}

	if pet.Category.Name != "" {
		var category models.Category
		db.DB.Where("name = ?", pet.Category.Name).FirstOrCreate(&category)
		pet.CategoryID = category.ID
	}

	var tags []models.Tag
	for _, tag := range pet.Tags {
		var existingTag models.Tag
		db.DB.Where("name = ?", tag.Name).FirstOrCreate(&existingTag)
		tags = append(tags, existingTag)
	}

	result := db.DB.Create(&pet)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to create pet"))
		return
	}

	if len(tags) > 0 {
		db.DB.Model(&pet).Association("Tags").Replace(tags)
	}

	c.JSON(http.StatusOK, pet)
}

func UpdatePet(c *gin.Context) {
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	var existingPet models.Pet
	if err := db.DB.First(&existingPet, pet.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Pet not found"))
		return
	}

	if pet.Category.Name != "" {
		var category models.Category
		db.DB.Where("name = ?", pet.Category.Name).FirstOrCreate(&category)
		pet.CategoryID = category.ID
	}

	var tags []models.Tag
	for _, tag := range pet.Tags {
		var existingTag models.Tag
		db.DB.Where("name = ?", tag.Name).FirstOrCreate(&existingTag)
		tags = append(tags, existingTag)
	}

	db.DB.Model(&existingPet).Updates(map[string]interface{}{
		"name":        pet.Name,
		"category_id": pet.CategoryID,
		"status":      pet.Status,
		"owner_id":    pet.OwnerID,
	})

	if len(tags) > 0 {
		db.DB.Model(&existingPet).Association("Tags").Replace(tags)
	}

	c.JSON(http.StatusOK, existingPet)
}

func FindPetsByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		status = "available" // Default as per spec
	}

	statuses := strings.Split(status, ",")
	var pets []models.Pet

	// Query pets by status
	db.DB.Where("status IN ?", statuses).Find(&pets)

	c.JSON(http.StatusOK, pets)
}

func FindPetsByTags(c *gin.Context) {
	tagsParam := c.QueryArray("tags")
	if len(tagsParam) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Tags parameter is required"))
		return
	}

	var pets []models.Pet
	db.DB.Joins("JOIN pet_tags ON pets.id = pet_tags.pet_id").
		Joins("JOIN tags ON pet_tags.tag_id = tags.id").
		Where("tags.name IN ?", tagsParam).
		Group("pets.id").
		Find(&pets)

	c.JSON(http.StatusOK, pets)
}

func GetPetById(c *gin.Context) {
	petId, err := strconv.ParseUint(c.Param("petId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid pet ID"))
		return
	}

	var pet models.Pet
	if err := db.DB.Preload("Category").Preload("Tags").First(&pet, petId).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Pet not found"))
		return
	}

	c.JSON(http.StatusOK, pet)
}

func UpdatePetWithForm(c *gin.Context) {
	petId, err := strconv.ParseUint(c.Param("petId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid pet ID"))
		return
	}

	name := c.Query("name")
	status := c.Query("status")

	var pet models.Pet
	if err := db.DB.First(&pet, petId).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Pet not found"))
		return
	}

	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if status != "" {
		updates["status"] = status
	}

	if len(updates) > 0 {
		db.DB.Model(&pet).Updates(updates)
	}

	c.JSON(http.StatusOK, pet)
}

func DeletePet(c *gin.Context) {
	petId, err := strconv.ParseUint(c.Param("petId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid pet ID"))
		return
	}

	var pet models.Pet
	if err := db.DB.First(&pet, petId).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Pet not found"))
		return
	}

	db.DB.Delete(&pet)

	c.JSON(http.StatusOK, gin.H{"message": "Pet deleted successfully"})
}

func UploadFile(c *gin.Context) {
	petId, err := strconv.ParseUint(c.Param("petId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid pet ID"))
		return
	}

	var pet models.Pet
	if err := db.DB.First(&pet, petId).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Pet not found"))
		return
	}

	additionalMetadata := c.Query("additionalMetadata")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "No file uploaded or invalid form"))
		return
	}
	defer file.Close()

	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}

	petDir := filepath.Join(uploadsDir, fmt.Sprintf("pet_%d", petId))
	if _, err := os.Stat(petDir); os.IsNotExist(err) {
		os.Mkdir(petDir, 0755)
	}

	filename := filepath.Join(petDir, header.Filename)

	out, err := os.Create(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to save file"))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to copy file"))
		return
	}

	response := models.ApiResponse{
		Code:    200,
		Type:    "success",
		Message: "File uploaded successfully. Additional metadata: " + additionalMetadata,
	}

	c.JSON(http.StatusOK, response)
}
