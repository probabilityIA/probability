package entities

type Recommendation struct {
	RecommendedCarrier string
	Reasoning          string
	Alternatives       []string
	Quotations         []Quotation
}

type Quotation struct {
	Carrier               string
	EstimatedCost         float64
	EstimatedDeliveryDays int
}
