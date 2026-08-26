package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"not null;index"`

	BusinessID    uint `gorm:"not null;index"`
	IntegrationID uint `gorm:"not null;default:0;index"`
	ConfigID      uint `gorm:"not null;default:0;index"`

	Module        string `gorm:"size:64;not null;default:'notification';index"`
	ReferenceType string `gorm:"size:64;not null;default:'';index:idx_email_logs_reference,priority:1"`
	ReferenceID   string `gorm:"size:64;not null;default:'';index:idx_email_logs_reference,priority:2"`

	To        string `gorm:"size:255;not null;index"`
	Subject   string `gorm:"size:512;not null"`
	EventType string `gorm:"size:128;not null;index"`

	Status            string  `gorm:"size:32;not null;index"`
	ErrorMessage      *string `gorm:"type:text"`
	Provider          string  `gorm:"size:32;not null;default:''"`
	ProviderMessageID string  `gorm:"size:128;not null;default:''"`

	SentBy     uint   `gorm:"not null;default:0"`
	SentByName string `gorm:"size:160;not null;default:''"`
}

func (EmailLog) TableName() string {
	return "email_logs"
}

func (e *EmailLog) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
