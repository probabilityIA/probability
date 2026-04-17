package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	WebhookSourceEnvioclick = "envioclick"

	WebhookLogStatusReceived  = "received"
	WebhookLogStatusProcessed = "processed"
	WebhookLogStatusFailed    = "failed"
	WebhookLogStatusIgnored   = "ignored"

	WebhookEventTrackingUpdate = "tracking_update"
	WebhookEventSync           = "sync"

	WebhookLogMaxPerSource = 50
)

type WebhookLog struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	Source         string
	EventType      string
	URL            string
	Method         string
	Headers        []byte
	RequestBody    []byte
	RemoteIP       string
	Status         string
	ResponseCode   int
	ProcessedAt    *time.Time
	ErrorMessage   *string
	ShipmentID     *uint
	BusinessID     *uint
	CorrelationID  *string
	TrackingNumber *string
	MappedStatus   *string
	RawStatus      *string
}
