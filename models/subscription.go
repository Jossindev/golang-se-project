package models

import (
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	gorm.Model
	Email     string    `json:"email" gorm:"uniqueIndex:idx_email_city" binding:"required,email"`
	City      string    `json:"city" gorm:"uniqueIndex:idx_email_city" binding:"required"`
	Frequency string    `json:"frequency" binding:"required,oneof=hourly daily"`
	Confirmed bool      `json:"confirmed" gorm:"default:false"`
	Token     string    `json:"token,omitempty" gorm:"index"`
	LastSent  time.Time `json:"-"`
}

type SubscriptionRequest struct {
	Email     string `form:"email" json:"email" binding:"required,email"`
	City      string `form:"city" json:"city" binding:"required"`
	Frequency string `form:"frequency" json:"frequency" binding:"required,oneof=hourly daily"`
}
