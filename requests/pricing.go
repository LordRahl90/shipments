package requests

type Pricing struct {
	Weight      float64 `json:"weight" form:"weight"`
	Origin      string  `json:"origin" form:"origin"`
	Destination string  `json:"destination" form:"destination"`
}
