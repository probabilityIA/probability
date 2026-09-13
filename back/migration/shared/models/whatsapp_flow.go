package models

import (
	"gorm.io/gorm"
)

type WhatsappFlow struct {
	gorm.Model

	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Name        string `gorm:"size:120;not null"`
	Description string `gorm:"type:text"`

	RootTemplateID *uint            `gorm:"index"`
	RootTemplate   WhatsappTemplate `gorm:"foreignKey:RootTemplateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Enabled bool `gorm:"not null;default:true"`
}

func (WhatsappFlow) TableName() string {
	return "whatsapp_flows"
}
