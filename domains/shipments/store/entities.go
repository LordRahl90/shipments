package store

import "gorm.io/gorm"

// Shipment db entity for shipment
type Shipment struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	CustomerID  string  `json:"customer_id"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Weight      float64 `json:"size"`
	Price       float64 `json:"price"`
	gorm.Model
}
