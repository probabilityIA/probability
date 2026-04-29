package models

import "gorm.io/gorm"

type ShippingMargin struct {
	gorm.Model
	BusinessID      uint    `gorm:"not null;index;uniqueIndex:idx_shipping_margin_biz_carrier,priority:1"`
	CarrierCode     string  `gorm:"size:50;not null;uniqueIndex:idx_shipping_margin_biz_carrier,priority:2"`
	CarrierName     string  `gorm:"size:100;not null"`
	MarginAmount    float64 `gorm:"type:decimal(15,2);not null;default:0"`
	InsuranceMargin float64 `gorm:"type:decimal(15,2);not null;default:0"`
	IsActive        bool    `gorm:"default:true;index"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
