package models

import "time"

type UserTourProgress struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null;index;uniqueIndex:idx_tour_user_business_key,priority:1"`
	BusinessID  uint       `gorm:"not null;default:0;index;uniqueIndex:idx_tour_user_business_key,priority:2"`
	TourKey     string     `gorm:"size:60;not null;uniqueIndex:idx_tour_user_business_key,priority:3"`
	Version     int        `gorm:"not null;default:1"`
	Status      string     `gorm:"size:20;not null;default:'pending';index"`
	StepIndex   int        `gorm:"not null;default:0"`
	CompletedAt *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserTourProgress) TableName() string {
	return "user_tour_progress"
}
