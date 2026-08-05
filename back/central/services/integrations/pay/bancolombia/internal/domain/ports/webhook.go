package ports

import (
	"context"
	"time"
)

type BancolombiaWebhookMessage struct {
	EventID        string     `json:"event_id"`
	Type           string     `json:"type"`
	TransferState  string     `json:"transfer_state"`
	Reference      string     `json:"reference"`
	TransferCode   string     `json:"transfer_code"`
	Amount         float64    `json:"amount"`
	Currency       string     `json:"currency"`
	PayerAccount   string     `json:"payer_account,omitempty"`
	OccurredAt     *time.Time `json:"occurred_at,omitempty"`
	IsTest         bool       `json:"is_test"`
	SignatureValid bool       `json:"signature_valid"`
	RawPayload     []byte     `json:"raw_payload"`
}

type IWebhookPublisher interface {
	PublishWebhookEvent(ctx context.Context, msg *BancolombiaWebhookMessage) error
}

type IWebhookUseCase interface {
	HandleIncomingWebhook(ctx context.Context, signatureHeader string, body []byte, isTest bool) error
}

type RawWebhookLog struct {
	ID              string
	Endpoint        string
	SignatureHeader string
	BodySize        int
	Body            []byte
	EventID         string
	EventType       string
	Reference       string
	TransferCode    string
}

type RawWebhookResult struct {
	ID          string
	Status      string
	HTTPStatus  int
	ErrorDetail string
}

type IRawWebhookLogger interface {
	LogIncoming(ctx context.Context, raw *RawWebhookLog) error
	UpdateResult(ctx context.Context, result *RawWebhookResult) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}
