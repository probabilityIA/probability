package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	SegmentCustomersInactive = "customers_inactive"
)

const (
	ScheduledRunStatusRunning   = "running"
	ScheduledRunStatusCompleted = "completed"
	ScheduledRunStatusFailed    = "failed"
)

const (
	ScheduledSendStatusQueued    = "queued"
	ScheduledSendStatusSent      = "sent"
	ScheduledSendStatusFailed    = "failed"
	ScheduledSendStatusDiscarded = "discarded"
)

type ScheduledNotificationRule struct {
	gorm.Model

	BusinessID uint     `gorm:"not null;index"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	IntegrationID *uint       `gorm:"index"`
	Integration   Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	NotificationTypeID uint             `gorm:"not null;index"`
	NotificationType   NotificationType `gorm:"foreignKey:NotificationTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	WhatsappTemplateID *uint            `gorm:"index"`
	WhatsappTemplate   WhatsappTemplate `gorm:"foreignKey:WhatsappTemplateID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Name        string `gorm:"size:160;not null"`
	Description string `gorm:"size:500"`

	SegmentType   string         `gorm:"size:64;not null;index"`
	SegmentParams datatypes.JSON `gorm:"type:jsonb"`

	Timezone         string `gorm:"size:64;not null;default:'America/Bogota'"`
	FrequencyMinutes uint   `gorm:"not null;default:1440"`
	SendWindowStart  string `gorm:"size:5;not null;default:'09:00'"`
	SendWindowEnd    string `gorm:"size:5;not null;default:'19:00'"`

	CooldownDays  uint `gorm:"not null;default:30"`
	DailySendCap  uint `gorm:"not null;default:500"`
	BatchSizeCap  uint `gorm:"not null;default:200"`
	RequiresOptIn bool `gorm:"not null;default:true"`

	Enabled   bool `gorm:"not null;default:true;index"`
	LastRunAt *time.Time
	NextRunAt *time.Time `gorm:"index"`

	CreatedByID *uint `gorm:"index"`
	CreatedBy   User  `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (ScheduledNotificationRule) TableName() string {
	return "scheduled_notification_rules"
}

type ScheduledNotificationRun struct {
	gorm.Model

	RuleID uint                      `gorm:"not null;index"`
	Rule   ScheduledNotificationRule `gorm:"foreignKey:RuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	BusinessID uint `gorm:"not null;index"`

	StartedAt  time.Time `gorm:"not null"`
	FinishedAt *time.Time

	Status string `gorm:"size:16;not null;default:'running';index"`

	MatchedCount   uint `gorm:"not null;default:0"`
	QueuedCount    uint `gorm:"not null;default:0"`
	SkippedCount   uint `gorm:"not null;default:0"`
	FailedCount    uint `gorm:"not null;default:0"`
	CapReachedFlag bool `gorm:"not null;default:false"`

	ErrorMessage string `gorm:"type:text"`
}

func (ScheduledNotificationRun) TableName() string {
	return "scheduled_notification_runs"
}

type ScheduledNotificationSend struct {
	gorm.Model

	RuleID uint                      `gorm:"not null;uniqueIndex:idx_sched_send_rule_client_day,priority:1"`
	Rule   ScheduledNotificationRule `gorm:"foreignKey:RuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ClientID uint   `gorm:"not null;uniqueIndex:idx_sched_send_rule_client_day,priority:2;index"`
	Client   Client `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	SendDate string `gorm:"size:10;not null;uniqueIndex:idx_sched_send_rule_client_day,priority:3"`

	RunID *uint `gorm:"index"`

	BusinessID uint   `gorm:"not null;index"`
	Phone      string `gorm:"size:32;not null"`

	Status       string `gorm:"size:16;not null;default:'queued';index"`
	MessageID    string `gorm:"size:128"`
	ErrorMessage string `gorm:"type:text"`

	QueuedAt time.Time `gorm:"not null;index"`
	SentAt   *time.Time
}

func (ScheduledNotificationSend) TableName() string {
	return "scheduled_notification_sends"
}
