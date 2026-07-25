package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	WebhookSourceWooCommerce = "woocommerce"

	WebhookLogStatusReceived = "received"
	WebhookLogStatusIgnored  = "ignored"
	WebhookLogStatusRejected = "rejected"

	WebhookLogRetentionDays = 30
)

type WebhookLog struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	EventType      string
	URL            string
	Method         string
	Headers        datatypes.JSON
	RequestBody    datatypes.JSON
	RemoteIP       string
	Status         string
	ResponseCode   int
	ErrorMessage   *string
	BusinessID     *uint
	IntegrationID  string
	RawStatus      *string
	MappedStatus   *string
	SignatureValid *bool
	RetentionUntil *time.Time
}

type IWebhookLogRepository interface {
	Save(ctx context.Context, entry *WebhookLog) error
	DeleteExpired(ctx context.Context) error
}
