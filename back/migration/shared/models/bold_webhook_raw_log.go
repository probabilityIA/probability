package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BoldWebhookRawLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index:idx_bold_webhook_raw_logs_created_at"`
	UpdatedAt time.Time `gorm:"not null"`

	Endpoint        string `gorm:"size:16;not null;index"`
	HTTPStatus      int    `gorm:"not null;default:0"`
	Status          string `gorm:"size:32;not null;index"`
	SignatureHeader string `gorm:"size:256"`

	BoldEventID       string `gorm:"size:128;index"`
	EventType         string `gorm:"size:64;index"`
	MerchantReference string `gorm:"size:128;index"`
	PaymentID         string `gorm:"size:128;index"`

	BodySize     int            `gorm:"not null;default:0"`
	BodyJSON     datatypes.JSON `gorm:"type:jsonb"`
	BodyText     *string        `gorm:"type:text"`
	ErrorDetail  *string        `gorm:"type:text"`
	ExpectedHash string         `gorm:"size:128"`
}

func (BoldWebhookRawLog) TableName() string {
	return "bold_webhook_raw_logs"
}

func (e *BoldWebhookRawLog) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
