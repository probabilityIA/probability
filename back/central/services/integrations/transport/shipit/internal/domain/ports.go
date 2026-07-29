package domain

import "context"

type IShipitClient interface {
	Rates(baseURL string, creds Credentials, req RatesRequest) (*RatesResponse, error)
	CreateShipment(baseURL string, creds Credentials, req ShipmentRequest) (*ShipmentResponse, error)
	Track(baseURL string, creds Credentials, trackingNumber string) (*TrackingResponse, error)
	GetCommunes(baseURL string, creds Credentials) ([]Commune, error)
	TestAuthentication(baseURL string, creds Credentials) error
}

type IWebhookResponsePublisher interface {
	PublishWebhookUpdate(ctx context.Context, msg *WebhookUpdateMessage) error
}

type WebhookUpdateMessage struct {
	ShipmentID      *uint
	BusinessID      uint
	CorrelationID   string
	TrackingNumber  string
	Carrier         string
	Status          ProbabilityShipmentStatus
	RawStatus       string
	RawStatusDetail string
	HasIncidence    bool
	IsUnknown       bool
	Description     string
	EventTimestamp  string
	ShippedAt       *string
	DeliveredAt     *string
}
