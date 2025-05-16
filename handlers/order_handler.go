package handlers

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetInventory(c *gin.Context) {
	// Count pets by status
	var statusCounts map[string]int64
	rows, err := db.DB.Model(&models.Pet{}).
		Select("status, count(*) as count").
		Group("status").
		Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to fetch inventory"))
		return
	}
	defer rows.Close()

	statusCounts = make(map[string]int64)

	for rows.Next() {
		var status string
		var count int64
		rows.Scan(&status, &count)
		statusCounts[status] = count
	}

	defaultStatuses := []string{"available", "pending", "sold"}
	for _, status := range defaultStatuses {
		if _, exists := statusCounts[status]; !exists {
			statusCounts[status] = 0
		}
	}

	c.JSON(http.StatusOK, statusCounts)
}

func PlaceOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	var pet models.Pet
	if err := db.DB.First(&pet, order.PetID).Error; err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Pet with ID "+strconv.FormatInt(order.PetID, 10)+" does not exist"))
		return
	}

	if order.Status == "" {
		order.Status = "placed"
	}

	now := time.Now()
	if order.ShipDate == nil {
		order.ShipDate = &now
	}

	if err := db.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to create order"))
		return
	}

	c.JSON(http.StatusOK, order.OrderResponse())
}

func GetOrderById(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid order ID"))
		return
	}

	if orderID > 10 && orderID <= 1000 {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Order not found"))
		return
	}

	var order models.Order
	if err := db.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Order not found"))
		return
	}

	c.JSON(http.StatusOK, order.OrderResponse())
}

func DeleteOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid order ID"))
		return
	}

	if orderID >= 1000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Order ID must be less than 1000"))
		return
	}

	var order models.Order
	if err := db.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Order not found"))
		return
	}

	if err := db.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to delete order"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
