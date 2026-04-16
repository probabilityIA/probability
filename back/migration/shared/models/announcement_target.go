package models

import "time"

type AnnouncementTarget struct {
	ID             uint      `gorm:"primaryKey"`
	AnnouncementID uint      `gorm:"not null;index;uniqueIndex:idx_announcement_business,priority:1"`
	BusinessID     uint      `gorm:"not null;index;uniqueIndex:idx_announcement_business,priority:2"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	Announcement Announcement `gorm:"foreignKey:AnnouncementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business     Business     `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (AnnouncementTarget) TableName() string {
	return "announcement_targets"
}
