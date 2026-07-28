package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PublicCheckout struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BusinessID      uint           `gorm:"not null;index"`
	IntegrationID   uint           `gorm:"not null;index"`
	Slug            string         `gorm:"size:100;not null;index"`
	Reference       string         `gorm:"size:100;not null;uniqueIndex"`
	Status          string         `gorm:"size:20;not null;default:'pending'"` // pending|paid|failed|expired
	Amount          float64        `gorm:"type:decimal(15,2);not null"`
	Currency        string         `gorm:"size:10;not null;default:'COP'"`
	Items           datatypes.JSON `gorm:"type:jsonb;not null"`
	CustomerName    string         `gorm:"size:255"`
	CustomerEmail   string         `gorm:"size:255"`
	CustomerPhone   string         `gorm:"size:50"`
	CustomerDni     string         `gorm:"size:50"`
	ShippingAddress datatypes.JSON `gorm:"type:jsonb"`
	BoldOrderID     string         `gorm:"size:100;index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PaidAt          *time.Time
}

func (PublicCheckout) TableName() string {
	return "public_checkouts"
}
