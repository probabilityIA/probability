package models

import (
	"gorm.io/gorm"
)

type WhatsappTemplateFlow struct {
	gorm.Model

	BusinessID uint     `gorm:"not null;index:idx_wa_flow_source_button,priority:1"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	FlowID *uint        `gorm:"index"`
	Flow   WhatsappFlow `gorm:"foreignKey:FlowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	SourceTemplateID uint             `gorm:"not null;index:idx_wa_flow_source_button,priority:2"`
	SourceTemplate   WhatsappTemplate `gorm:"foreignKey:SourceTemplateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ButtonText string `gorm:"size:64;not null;index:idx_wa_flow_source_button,priority:3"`

	TargetTemplateID uint             `gorm:"not null;index"`
	TargetTemplate   WhatsappTemplate `gorm:"foreignKey:TargetTemplateID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Enabled bool `gorm:"not null;default:true"`
}

func (WhatsappTemplateFlow) TableName() string {
	return "whatsapp_template_flows"
}
