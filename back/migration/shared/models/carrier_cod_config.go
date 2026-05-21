package models

import "gorm.io/gorm"

type CarrierCodConfig struct {
	gorm.Model
	BusinessID         uint    `gorm:"not null;index;uniqueIndex:idx_carrier_cod_biz_carrier,priority:1"`
	CarrierName        string  `gorm:"size:128;not null;uniqueIndex:idx_carrier_cod_biz_carrier,priority:2"`
	DiscountPercentage float64 `gorm:"type:decimal(6,2);not null;default:0"`
	IsActive           bool    `gorm:"not null;default:true;index"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
