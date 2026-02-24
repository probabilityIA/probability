package domain

import "time"

// EnvioClick API Models - copied from central/services/integrations/transport/envioclick/internal/domain/entities.go

type QuoteRequest struct {
	IDRate              int64     `json:"idRate"`
	MyShipmentReference string    `json:"myShipmentReference"`
	ExternalOrderID     string    `json:"external_order_id"`
	OrderUUID           string    `json:"order_uuid,omitempty"`
	RequestPickup       bool      `json:"requestPickup"`
	PickupDate          string    `json:"pickupDate"`
	Insurance           bool      `json:"insurance"`
	Description         string    `json:"description"`
	ContentValue        float64   `json:"contentValue"`
	CODValue            float64   `json:"codValue"`
	IncludeGuideCost    bool      `json:"includeGuideCost"`
	CODPaymentMethod    string    `json:"codPaymentMethod"`
	Packages            []Package `json:"packages"`
	Origin              Address   `json:"origin"`
	Destination         Address   `json:"destination"`
}

type Package struct {
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
}

type Address struct {
	Company     string `json:"company"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Suburb      string `json:"suburb"`
	CrossStreet string `json:"crossStreet"`
	Reference   string `json:"reference"`
	DaneCode    string `json:"daneCode"`
}

type QuoteResponse struct {
	Status string    `json:"status"`
	Data   QuoteData `json:"data"`
}

type QuoteData struct {
	Rates []Rate `json:"rates"`
}

type Rate struct {
	IDRate        int64   `json:"idRate"`
	IDProduct     int64   `json:"idProduct"`
	Product       string  `json:"product"`
	IDCarrier     int64   `json:"idCarrier"`
	Carrier       string  `json:"carrier"`
	Flete         float64 `json:"flete"`
	DeliveryDays  int     `json:"deliveryDays"`
	QuotationType string  `json:"quotationType"`
}

type GenerateResponse struct {
	Status string       `json:"status"`
	Data   GenerateData `json:"data"`
}

type GenerateData struct {
	TrackingNumber   string `json:"tracker"`
	LabelURL         string `json:"url"`
	MyGuideReference string `json:"myGuideReference"`
}

type TrackingResponse struct {
	Status string       `json:"status"`
	Data   TrackingData `json:"data"`
}

type TrackingData struct {
	TrackingNumber string         `json:"trackingNumber"`
	Carrier        string         `json:"carrier"`
	Status         string         `json:"status"`
	Events         []TrackHistory `json:"history"`
}

type TrackHistory struct {
	Date        string `json:"date"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type CancelResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// TrackRequest is the body sent by the real client to /api/v2/track
type TrackRequest struct {
	TrackingCode string `json:"trackingCode"`
}

// StoredShipment is the internal model for in-memory storage
type StoredShipment struct {
	ID             string
	TrackingNumber string
	Carrier        string
	CarrierID      int64
	Product        string
	Origin         Address
	Destination    Address
	Packages       []Package
	ContentValue   float64
	Flete          float64
	Status         string // created, in_transit, delivered, cancelled
	LabelURL       string
	CreatedAt      time.Time
	CancelledAt    *time.Time
}
