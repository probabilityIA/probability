package models

import "time"

type AnnouncementView struct {
	ID             uint      `gorm:"primaryKey"`
	AnnouncementID uint      `gorm:"not null;index;index:idx_announcement_user,priority:1"`
	UserID         uint      `gorm:"not null;index;index:idx_announcement_user,priority:2"`
	BusinessID     uint      `gorm:"not null;index"`
	Action         string    `gorm:"size:20;not null;index"`
	LinkID         *uint     `gorm:"index"`
	ViewedAt       time.Time `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	Announcement     Announcement      `gorm:"foreignKey:AnnouncementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User             User              `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business         Business          `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AnnouncementLink *AnnouncementLink `gorm:"foreignKey:LinkID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (AnnouncementView) TableName() string {
	return "announcement_views"
}
