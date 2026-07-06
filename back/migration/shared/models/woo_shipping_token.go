package models

import "gorm.io/gorm"

type WooShippingToken struct {
	gorm.Model
	IntegrationID uint   `gorm:"not null;uniqueIndex"`
	Salt          string `gorm:"size:64;not null"`
	Revoked       bool   `gorm:"not null;default:false"`
}

func (WooShippingToken) TableName() string {
	return "woo_shipping_tokens"
}
