package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SubscriptionType struct {
	gorm.Model
	Name                 string         `gorm:"size:100;not null"`
	Code                 string         `gorm:"size:50;not null;unique"`
	Description          string         `gorm:"type:text"`
	Price                float64        `gorm:"type:numeric;not null"`
	BillingPeriod        string         `gorm:"size:20;not null;default:'monthly'"`
	Active               bool           `gorm:"default:true"`
	Features             datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'"`
	MaxEcommerceChannels int            `gorm:"not null;default:0"`
}

func (SubscriptionType) TableName() string {
	return "subscription_types"
}
