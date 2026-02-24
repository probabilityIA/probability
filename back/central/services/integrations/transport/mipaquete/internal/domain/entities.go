package domain

// ═══════════════════════════════════════════════════════════════
// MiPaquete transport DTOs
// TODO: Replace with actual MiPaquete API request/response models
// ═══════════════════════════════════════════════════════════════

type QuoteRequest struct {
	Origin      Location  `json:"origin"`
	Destination Location  `json:"destination"`
	Packages    []Package `json:"packages"`
}

type Location struct {
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	ZipCode    string `json:"zip_code,omitempty"`
	LocationID string `json:"location_id,omitempty"`
}

type Package struct {
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type QuoteResponse struct {
	Rates []Rate `json:"rates"`
}

type Rate struct {
	Carrier      string  `json:"carrier"`
	Service      string  `json:"service"`
	Price        float64 `json:"price"`
	DeliveryDays int     `json:"delivery_days"`
}

type GenerateResponse struct {
	ShipmentID     string `json:"shipment_id"`
	TrackingNumber string `json:"tracking_number"`
	LabelURL       string `json:"label_url,omitempty"`
	Carrier        string `json:"carrier"`
}

type TrackingResponse struct {
	TrackingNumber string          `json:"tracking_number"`
	Status         string          `json:"status"`
	Events         []TrackingEvent `json:"events"`
}

type TrackingEvent struct {
	Date        string `json:"date"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location,omitempty"`
}

type CancelResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
