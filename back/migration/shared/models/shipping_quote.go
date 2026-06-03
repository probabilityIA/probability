package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShippingQuote struct {
	gorm.Model

	BusinessID    uint   `gorm:"not null;index"`
	IntegrationID uint   `gorm:"index"`
	Source        string `gorm:"size:32;not null;index"`
	CorrelationID string `gorm:"size:64;index"`

	OrderUUID        *string `gorm:"size:64;index"`
	ExternalOrderRef string  `gorm:"size:128"`

	RequestPayload datatypes.JSON `gorm:"type:jsonb"`
	Rates          datatypes.JSON `gorm:"type:jsonb"`

	SelectedCarrier     string `gorm:"size:128"`
	SelectedServiceCode string `gorm:"size:128"`
	SelectedIDRate      *int64

	Status    string     `gorm:"size:32;not null;default:'created';index"`
	ExpiresAt *time.Time `gorm:"index"`
}

func (ShippingQuote) TableName() string {
	return "shipping_quotes"
}
