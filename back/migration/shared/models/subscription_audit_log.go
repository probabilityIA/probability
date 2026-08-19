package models

import "time"

type SubscriptionAuditLog struct {
	ID          uint      `gorm:"primaryKey"`
	BusinessID  uint      `gorm:"not null;index"`
	ActorUserID *uint     `gorm:"index"`
	ActorLabel  string    `gorm:"size:150;not null"`
	Action      string    `gorm:"size:60;not null;index"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (SubscriptionAuditLog) TableName() string {
	return "subscription_audit_logs"
}
