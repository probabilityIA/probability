package dtos

import (
	"time"

	"github.com/google/uuid"
)

const (
	GatewayBold        = "bold"
	GatewayBancolombia = "bancolombia"
)

type PaymentWebhookEvent struct {
	ID                   uuid.UUID
	Gateway              string
	EventID              string
	Type                 string
	ExternalID           string
	ExternalState        string
	Reference            string
	Source               string
	OccurredAt           *time.Time
	Payload              []byte
	SignatureValid       bool
	ProcessedAt          *time.Time
	ProcessedError       *string
	PaymentTransactionID *uint
}
