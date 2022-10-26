package entities

import "shipments/domains/shipments/store"

// Shipment contains a shipment service entity
type Shipment struct {
	ID          string  `json:"id"`
	CustomerID  string  `json:"customer_id"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Weight      float64 `json:"float"`
	Price       float64 `json:"price"`
}

// ToDBShipment converts shipment service entity to db entity
func (s *Shipment) ToDBShipment() *store.Shipment {
	return &store.Shipment{
		ID:          s.ID,
		CustomerID:  s.CustomerID,
		Origin:      s.Origin,
		Destination: s.Destination,
		Weight:      s.Weight,
		Price:       s.Price,
	}
}

// FromDBShipmentEntity converts from db entity to service entity
func FromDBShipmentEntity(s *store.Shipment) *Shipment {
	return &Shipment{
		ID:          s.ID,
		CustomerID:  s.CustomerID,
		Origin:      s.Origin,
		Destination: s.Destination,
		Weight:      s.Weight,
		Price:       s.Price,
	}
}
