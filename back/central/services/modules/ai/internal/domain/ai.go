package domain

type RecommendationRequest struct {
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	Carriers    []string `json:"carriers,omitempty"` // Carriers from quotes
}

type AIRecommendation struct {
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

type AIService interface {
	GetRecommendation(req RecommendationRequest) (*AIRecommendation, error)
}
