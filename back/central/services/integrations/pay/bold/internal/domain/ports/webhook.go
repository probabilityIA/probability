package ports

import (
	"context"
	"time"
)

type BoldWebhookMessage struct {
	BoldEventID       string     `json:"bold_event_id"`
	Type              string     `json:"type"`
	Subject           string     `json:"subject"`
	Source            string     `json:"source"`
	OccurredAt        *time.Time `json:"occurred_at,omitempty"`
	PaymentID         string     `json:"payment_id"`
	MerchantReference string     `json:"merchant_reference"`
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	PaymentMethod     string     `json:"payment_method"`
	PayerEmail        string     `json:"payer_email,omitempty"`
	IsTest            bool       `json:"is_test"`
	RawPayload        []byte     `json:"raw_payload"`
}

type IWebhookPublisher interface {
	PublishWebhookEvent(ctx context.Context, msg *BoldWebhookMessage) error
}

type IWebhookUseCase interface {
	HandleIncomingWebhook(ctx context.Context, signatureHeader string, body []byte, isTest bool) error
}

type RawWebhookLog struct {
	ID                string
	Endpoint          string
	SignatureHeader   string
	BodySize          int
	Body              []byte
	BoldEventID       string
	EventType         string
	MerchantReference string
	PaymentID         string
}

type RawWebhookResult struct {
	ID           string
	Status       string
	HTTPStatus   int
	ErrorDetail  string
	ExpectedHash string
}

type IRawWebhookLogger interface {
	LogIncoming(ctx context.Context, raw *RawWebhookLog) error
	UpdateResult(ctx context.Context, result *RawWebhookResult) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}
