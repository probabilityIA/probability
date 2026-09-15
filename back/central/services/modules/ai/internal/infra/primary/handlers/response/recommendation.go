package response

type Recommendation struct {
	RecommendedCarrier string      `json:"recommended_carrier"`
	Reasoning          string      `json:"reasoning"`
	Alternatives       []string    `json:"alternatives"`
	Quotations         []Quotation `json:"quotations"`
}

type Quotation struct {
	Carrier               string  `json:"carrier"`
	EstimatedCost         float64 `json:"estimated_cost"`
	EstimatedDeliveryDays int     `json:"estimated_delivery_days"`
}
