package responses

// Pricing response for pricing requests
type Pricing struct {
	Weight         float64 `json:"weight"`
	Origin         string  `json:"origin"`
	Destination    string  `json:"destination"`
	WeightCategory string  `json:"weight_category"`
	Price          float64 `json:"price"`
}
