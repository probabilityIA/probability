package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductFieldOrigin struct {
	gorm.Model

	ProductID  string `gorm:"type:varchar(64);not null;uniqueIndex:idx_field_origin_key,priority:1"`
	Field      string `gorm:"size:40;not null;uniqueIndex:idx_field_origin_key,priority:2"`
	BusinessID uint   `gorm:"not null;index"`

	Source        string `gorm:"size:16;not null;index"`
	IntegrationID *uint  `gorm:"index"`
	UserID        *uint  `gorm:"index"`

	ChangedAt time.Time `gorm:"not null;index"`

	Product     Product      `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business    Business     `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integration *Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (ProductFieldOrigin) TableName() string {
	return "product_field_origins"
}

type ProductFieldChange struct {
	ID uint `gorm:"primarykey"`

	ProductID  string `gorm:"type:varchar(64);not null;index:idx_field_change_key,priority:1"`
	Field      string `gorm:"size:40;not null;index:idx_field_change_key,priority:2"`
	BusinessID uint   `gorm:"not null;index"`

	OldValue  string `gorm:"type:text"`
	Truncated bool   `gorm:"default:false"`

	Source        string `gorm:"size:16;not null;index"`
	IntegrationID *uint  `gorm:"index"`
	UserID        *uint  `gorm:"index"`

	BatchID string `gorm:"size:64;index"`

	CreatedAt time.Time `gorm:"not null;index:idx_field_change_key,priority:3,sort:desc"`

	Product  Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ProductFieldChange) TableName() string {
	return "product_field_changes"
}
