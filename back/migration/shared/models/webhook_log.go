package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WebhookLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index:idx_webhook_logs_source_created,priority:2"`
	UpdatedAt time.Time `gorm:"not null"`

	Source    string `gorm:"size:64;not null;index:idx_webhook_logs_source_created,priority:1"`
	EventType string `gorm:"size:128;not null;index"`

	URL         string         `gorm:"size:512;not null"`
	Method      string         `gorm:"size:8;not null;default:POST"`
	Headers     datatypes.JSON `gorm:"type:jsonb"`
	RequestBody datatypes.JSON `gorm:"type:jsonb;not null"`
	RemoteIP    string         `gorm:"size:64"`

	Status       string     `gorm:"size:32;not null;index"`
	ResponseCode int        `gorm:"not null;default:200"`
	ProcessedAt  *time.Time
	ErrorMessage *string `gorm:"type:text"`

	ShipmentID    *uint   `gorm:"index"`
	BusinessID    *uint   `gorm:"index"`
	CorrelationID *string `gorm:"size:128;index"`

	TrackingNumber *string `gorm:"size:128;index"`
	MappedStatus   *string `gorm:"size:64"`
	RawStatus      *string `gorm:"size:128"`
}

func (WebhookLog) TableName() string {
	return "webhook_logs"
}

func (w *WebhookLog) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
