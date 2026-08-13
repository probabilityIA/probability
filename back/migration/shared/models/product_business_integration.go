package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductBusinessIntegration struct {
	gorm.Model

	ProductID string `gorm:"type:varchar(64);not null;index"`

	BusinessID uint `gorm:"not null;index"`

	IntegrationID uint `gorm:"not null;index"`

	ExternalProductID string  `gorm:"size:255;not null;index"`
	ExternalVariantID *string `gorm:"size:255;index"`
	ExternalSKU       *string `gorm:"size:255;index"`
	ExternalBarcode   *string `gorm:"size:255;index"`

	ExternalLogisticType *string `gorm:"size:32;index"`

	LastPushedQty *int `gorm:"column:last_pushed_qty"`

	ChannelSnapshot datatypes.JSON `gorm:"column:channel_snapshot;type:jsonb"`
	SnapshotAt      *time.Time     `gorm:"column:snapshot_at;index"`

	Product     Product     `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business    Business    `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integration Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ProductBusinessIntegration) TableName() string {
	return "product_business_integrations"
}
