package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerificationToken struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	TokenHash string `gorm:"size:64;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null;index"`
	UsedAt    *time.Time
	User      *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}
