package store

import (
	"gorm.io/gorm"
)

// Customer model representing the customer entity
type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"size:50;uniqueIndex"`
	gorm.Model
}
