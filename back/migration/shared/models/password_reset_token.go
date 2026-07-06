package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"size:64;not null;uniqueIndex"`
	Channel   string    `gorm:"size:20;not null;default:email"`
	Attempts  int       `gorm:"not null;default:0"`
	ExpiresAt time.Time `gorm:"not null;index"`
	UsedAt    *time.Time
	User      *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
