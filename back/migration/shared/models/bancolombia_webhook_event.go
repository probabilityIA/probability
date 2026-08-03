package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BancolombiaWebhookEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index"`
	UpdatedAt time.Time `gorm:"not null"`

	EventID       string `gorm:"size:128;not null;uniqueIndex:uq_bancolombia_webhook_event_id"`
	Type          string `gorm:"size:64;not null;index"`
	TransferState string `gorm:"size:64;index"`
	Reference     string `gorm:"size:128;index"`
	OccurredAt    *time.Time

	Payload datatypes.JSON `gorm:"type:jsonb;not null"`

	SignatureValid bool `gorm:"not null;default:false;index"`

	ProcessedAt    *time.Time `gorm:"index"`
	ProcessedError *string    `gorm:"type:text"`

	PaymentTransactionID *uint `gorm:"index"`
}

func (BancolombiaWebhookEvent) TableName() string {
	return "bancolombia_webhook_events"
}

func (e *BancolombiaWebhookEvent) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
