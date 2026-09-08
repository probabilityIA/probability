package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	CampaignStatusDraft     = "draft"
	CampaignStatusScheduled = "scheduled"
	CampaignStatusRunning   = "running"
	CampaignStatusPaused    = "paused"
	CampaignStatusCompleted = "completed"
	CampaignStatusCancelled = "cancelled"
)

const (
	CampaignSendStatusPending   = "pending"
	CampaignSendStatusQueued    = "queued"
	CampaignSendStatusSent      = "sent"
	CampaignSendStatusDelivered = "delivered"
	CampaignSendStatusRead      = "read"
	CampaignSendStatusReplied   = "replied"
	CampaignSendStatusFailed    = "failed"
	CampaignSendStatusSkipped   = "skipped"
)

const (
	CampaignAudienceAllClients = "all_clients"
	CampaignAudienceFiltered   = "filtered_clients"
)

type WhatsappCampaign struct {
	gorm.Model

	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	IntegrationID *uint       `gorm:"index"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	WhatsappTemplateID *uint            `gorm:"index"`
	WhatsappTemplate   WhatsappTemplate `gorm:"foreignKey:WhatsappTemplateID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Name        string `gorm:"size:160;not null"`
	Description string `gorm:"size:500"`

	SenderName string `gorm:"size:120"`

	AudienceType   string         `gorm:"size:32;not null;default:'filtered_clients'"`
	AudienceParams datatypes.JSON `gorm:"type:jsonb"`
	VariableValues datatypes.JSON `gorm:"type:jsonb"`

	Timezone        string `gorm:"size:64;not null;default:'America/Bogota'"`
	SendWindowStart string `gorm:"size:5;not null;default:'09:00'"`
	SendWindowEnd   string `gorm:"size:5;not null;default:'19:00'"`

	ScheduledAt  *time.Time `gorm:"index"`
	DailySendCap uint       `gorm:"not null;default:250"`
	BatchSize    uint       `gorm:"not null;default:50"`

	Status string `gorm:"size:16;not null;default:'draft';index"`

	AudienceCount uint `gorm:"not null;default:0"`
	QueuedCount   uint `gorm:"not null;default:0"`
	SentCount     uint `gorm:"not null;default:0"`
	FailedCount   uint `gorm:"not null;default:0"`
	SkippedCount  uint `gorm:"not null;default:0"`
	RepliedCount  uint `gorm:"not null;default:0"`

	StartedAt   *time.Time
	FinishedAt  *time.Time
	LastBatchAt *time.Time `gorm:"index"`

	ErrorMessage string `gorm:"type:text"`

	CreatedByID *uint `gorm:"index"`
	CreatedBy   User  `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (WhatsappCampaign) TableName() string {
	return "whatsapp_campaigns"
}

type WhatsappCampaignSend struct {
	gorm.Model

	CampaignID uint             `gorm:"not null;uniqueIndex:idx_campaign_send_campaign_client,priority:1"`
	Campaign   WhatsappCampaign `gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ClientID uint   `gorm:"not null;uniqueIndex:idx_campaign_send_campaign_client,priority:2;index"`
	Client   Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BusinessID uint   `gorm:"not null;index"`
	Phone      string `gorm:"size:32;not null;index"`

	ConversationID *uuid.UUID `gorm:"type:uuid;index"`

	Status       string `gorm:"size:16;not null;default:'pending';index"`
	MessageID    string `gorm:"size:128;index"`
	ErrorMessage string `gorm:"type:text"`

	QueuedAt    *time.Time
	SentAt      *time.Time
	DeliveredAt *time.Time
	ReadAt      *time.Time
	RepliedAt   *time.Time
}

func (WhatsappCampaignSend) TableName() string {
	return "whatsapp_campaign_sends"
}
