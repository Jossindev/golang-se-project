package models

import (
	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	Name       string   `json:"name" binding:"required"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	CategoryID uint     `json:"-"`
	PhotoUrls  []string `json:"photoUrls" gorm:"-"`
	Tags       []Tag    `json:"tags" gorm:"many2many:pet_tags;"`
	Status     string   `json:"status" binding:"omitempty,oneof=available pending sold"`
	OwnerID    uint     `json:"ownerId"`
}
