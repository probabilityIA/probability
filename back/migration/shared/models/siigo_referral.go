package models

import "gorm.io/gorm"

type SiigoReferral struct {
	gorm.Model
	Name       string `gorm:"size:255;not null"`
	Email      string `gorm:"size:255;not null"`
	Phone      string `gorm:"size:50;not null"`
	OrderRange string `gorm:"size:20;not null"`
}

func (SiigoReferral) TableName() string {
	return "siigo_referrals"
}
