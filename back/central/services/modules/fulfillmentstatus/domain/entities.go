package domain

// FulfillmentStatusInfo contiene información básica del estado de fulfillment
type FulfillmentStatusInfo struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Color       string `json:"color"`
}
