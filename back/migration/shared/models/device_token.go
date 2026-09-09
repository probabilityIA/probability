package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceToken struct {
	gorm.Model

	UserID     uint  `gorm:"not null;index"`
	BusinessID *uint `gorm:"index"`

	Token    string `gorm:"size:512;not null;uniqueIndex"`
	Platform string `gorm:"size:20;not null;index"`

	AppVersion string `gorm:"size:50"`
	DeviceName string `gorm:"size:120"`

	IsActive   bool `gorm:"default:true;index"`
	LastSeenAt *time.Time

	User     User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (DeviceToken) TableName() string {
	return "device_tokens"
}
