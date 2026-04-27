package dtos

import (
	"time"

	"github.com/google/uuid"
)

type BoldWebhookEvent struct {
	ID                   uuid.UUID
	BoldEventID          string
	Type                 string
	Subject              string
	Source               string
	OccurredAt           *time.Time
	Payload              []byte
	SignatureValid       bool
	ProcessedAt          *time.Time
	ProcessedError       *string
	PaymentTransactionID *uint
}

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
	RawPayload        []byte     `json:"raw_payload"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
}
