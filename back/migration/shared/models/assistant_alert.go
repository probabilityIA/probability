package models

import "time"

type AssistantAlert struct {
	ID string `gorm:"type:uuid;primaryKey"`

	BusinessID uint      `gorm:"not null;index:idx_assistant_alerts_business_created,priority:1"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	EventID   string `gorm:"size:80;not null;uniqueIndex"`
	EventType string `gorm:"size:80;not null;index"`
	Severity  string `gorm:"size:20;not null;default:'info'"`

	Title string `gorm:"size:160;not null"`
	Body  string `gorm:"type:text;not null"`

	DestinationKey   string `gorm:"size:80"`
	DestinationRoute string `gorm:"size:255"`
	ReferenceType    string `gorm:"size:40"`
	ReferenceID      string `gorm:"size:80"`

	CreatedAt time.Time `gorm:"not null;index:idx_assistant_alerts_business_created,priority:2,sort:desc"`
}

func (AssistantAlert) TableName() string {
	return "assistant_alerts"
}

type AssistantAlertCursor struct {
	UserID     uint `gorm:"primaryKey;autoIncrement:false"`
	BusinessID uint `gorm:"primaryKey;autoIncrement:false"`

	SeenAt    time.Time `gorm:"not null"`
	UpdatedAt time.Time
}

func (AssistantAlertCursor) TableName() string {
	return "assistant_alert_cursors"
}
