package dtos

import (
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
	PublishedAt    *time.Time `json:"published_at,omitempty"`
}
