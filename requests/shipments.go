package requests

// Shipment content of a shipment request
type Shipment struct {
	Name        string  `json:"name" validate:"required"`
	Email       string  `json:"email"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Weight      float64 `json:"weight"`
}
