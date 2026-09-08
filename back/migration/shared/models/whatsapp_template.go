package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	WhatsappTemplateStatusDraft    = "draft"
	WhatsappTemplateStatusPending  = "pending"
	WhatsappTemplateStatusApproved = "approved"
	WhatsappTemplateStatusRejected = "rejected"
	WhatsappTemplateStatusPaused   = "paused"
	WhatsappTemplateStatusDisabled = "disabled"
	WhatsappTemplateStatusFailed   = "failed"
)

const (
	WhatsappTemplateCategoryMarketing = "MARKETING"
	WhatsappTemplateCategoryUtility   = "UTILITY"
	WhatsappTemplateCategoryAuth      = "AUTHENTICATION"
)

const (
	WhatsappTemplateOriginSystem   = "system"
	WhatsappTemplateOriginInternal = "internal"
	WhatsappTemplateOriginBusiness = "business"
)

const (
	WhatsappTemplateScopeOrderEvent = "order_event"
	WhatsappTemplateScopeScheduled  = "scheduled"
	WhatsappTemplateScopeInternal   = "internal"
)

type WhatsappTemplate struct {
	gorm.Model

	BusinessID *uint    `gorm:"index:idx_whatsapp_templates_business_name,priority:1"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Origin string `gorm:"size:16;not null;default:'business';index"`
	Scope  string `gorm:"size:24;not null;default:'scheduled';index"`

	IntegrationID *uint       `gorm:"index"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Name     string `gorm:"size:512;not null;index:idx_whatsapp_templates_business_name,priority:2"`
	Language string `gorm:"size:16;not null;default:'es';index:idx_whatsapp_templates_business_name,priority:3"`
	Category string `gorm:"size:32;not null;default:'MARKETING'"`

	BodyText   string `gorm:"type:text;not null"`
	HeaderText string `gorm:"type:text"`
	FooterText string `gorm:"type:text"`

	VariableMapping datatypes.JSON `gorm:"type:jsonb"`
	Buttons         datatypes.JSON `gorm:"type:jsonb"`
	Components      datatypes.JSON `gorm:"type:jsonb"`

	WABAID         string `gorm:"size:64;index"`
	MetaTemplateID string `gorm:"size:64;index"`

	Status         string `gorm:"size:24;not null;default:'draft';index"`
	RejectedReason string `gorm:"type:text"`
	SubmittedAt    *time.Time
	ReviewedAt     *time.Time
	LastSyncedAt   *time.Time

	CreatedByID *uint `gorm:"index"`
	CreatedBy   User  `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (WhatsappTemplate) TableName() string {
	return "whatsapp_templates"
}
