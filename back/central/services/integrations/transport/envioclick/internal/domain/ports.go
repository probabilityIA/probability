package domain

import (
	"context"

	"github.com/google/uuid"
)

type IEnvioClickClient interface {
	Quote(baseURL, apiKey string, req QuoteRequest) (*QuoteResponse, error)
	Generate(baseURL, apiKey string, req QuoteRequest) (*GenerateResponse, error)
	Track(baseURL, apiKey string, trackingNumber string) (*TrackingResponse, error)
	TrackByOrdersBatch(baseURL, apiKey string, orders []int64) (*TrackingResponse, error)
	Cancel(baseURL, apiKey string, idShipment string) (*CancelResponse, error)
	CancelBatch(baseURL, apiKey string, req CancelBatchRequest) (*CancelBatchResponse, error)
}

type IWebhookLogRepository interface {
	Save(ctx context.Context, log *WebhookLog) error
	MarkProcessed(ctx context.Context, id uuid.UUID, errorMessage *string) error
	TrimOldBySource(ctx context.Context, source string, keepCount int) error
}

type IWebhookResponsePublisher interface {
	PublishWebhookUpdate(ctx context.Context, msg *WebhookUpdateMessage) error
}

type WebhookUpdateMessage struct {
	ShipmentID     *uint
	BusinessID     uint
	CorrelationID  string
	TrackingNumber string
	Status         ProbabilityShipmentStatus
	RawStatus      string
	HasIncidence   bool
	IsUnknown      bool
	Description    string
	EventTimestamp string
	ShippedAt      *string
	DeliveredAt    *string
}
