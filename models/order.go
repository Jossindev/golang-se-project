package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	PetID    int64      `json:"petId" binding:"required"`
	Quantity int32      `json:"quantity" binding:"required"`
	ShipDate *time.Time `json:"shipDate,omitempty"`
	Status   string     `json:"status,omitempty" binding:"omitempty,oneof=placed approved delivered canceled"`
	Complete bool       `json:"complete,omitempty"`
}

func (o *Order) OrderResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":        o.ID,
		"petId":     o.PetID,
		"quantity":  o.Quantity,
		"shipDate":  o.ShipDate,
		"status":    o.Status,
		"complete":  o.Complete,
		"createdAt": o.CreatedAt,
	}
}
