package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BoldWebhookEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index"`
	UpdatedAt time.Time `gorm:"not null"`

	BoldEventID string `gorm:"size:128;not null;uniqueIndex:uq_bold_webhook_event_id"`
	Type        string `gorm:"size:64;not null;index"`
	Subject     string `gorm:"size:128;index"`
	Source      string `gorm:"size:128"`
	OccurredAt  *time.Time

	Payload datatypes.JSON `gorm:"type:jsonb;not null"`

	SignatureValid bool `gorm:"not null;default:false;index"`

	ProcessedAt    *time.Time `gorm:"index"`
	ProcessedError *string    `gorm:"type:text"`

	PaymentTransactionID *uint      `gorm:"index"`
	WalletTransactionID  *uuid.UUID `gorm:"type:uuid;index"`
}

func (BoldWebhookEvent) TableName() string {
	return "bold_webhook_events"
}

func (e *BoldWebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
