package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string `json:"username" gorm:"unique"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	UserStatus int    `json:"userStatus,omitempty"` // User status
}
