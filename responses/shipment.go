package responses

type Shipment struct {
	Success    bool    `json:"success"`
	Reference  string  `json:"reference"`
	CustomerID string  `json:"customer_id"`
	Price      float64 `json:"price"`
}
