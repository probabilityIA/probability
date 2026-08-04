package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PaymentWebhookEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index"`
	UpdatedAt time.Time `gorm:"not null"`

	Gateway string `gorm:"size:32;not null;uniqueIndex:uq_payment_webhook_gateway_event,priority:1"`
	EventID string `gorm:"size:128;not null;uniqueIndex:uq_payment_webhook_gateway_event,priority:2"`

	Type          string `gorm:"size:64;not null;index"`
	ExternalID    string `gorm:"size:128;index"`
	ExternalState string `gorm:"size:64;index"`
	Reference     string `gorm:"size:128;index"`
	Source        string `gorm:"size:128"`
	OccurredAt    *time.Time

	Payload datatypes.JSON `gorm:"type:jsonb;not null"`

	SignatureValid bool `gorm:"not null;default:false;index"`

	ProcessedAt    *time.Time `gorm:"index"`
	ProcessedError *string    `gorm:"type:text"`

	PaymentTransactionID *uint      `gorm:"index"`
	WalletTransactionID  *uuid.UUID `gorm:"type:uuid;index"`
}

func (PaymentWebhookEvent) TableName() string {
	return "payment_webhook_events"
}

func (e *PaymentWebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
