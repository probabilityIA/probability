package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CredentialRevealAudit struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt         time.Time `gorm:"not null;index"`
	UpdatedAt         time.Time `gorm:"not null"`
	UserID            uint      `gorm:"not null;index"`
	BusinessID        uint      `gorm:"index"`
	IntegrationTypeID uint      `gorm:"not null;index"`
	IntegrationCode   string    `gorm:"size:64;not null;index"`
	IPAddress         string    `gorm:"size:64"`
	UserAgent         string    `gorm:"size:512"`
}

func (CredentialRevealAudit) TableName() string {
	return "credential_reveal_audits"
}

func (a *CredentialRevealAudit) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
