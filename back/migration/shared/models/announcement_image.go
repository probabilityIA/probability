package models

import "time"

type AnnouncementImage struct {
	ID             uint      `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	AnnouncementID uint      `gorm:"not null;index"`
	ImageURL       string    `gorm:"size:500;not null"`
	SortOrder      int       `gorm:"default:0"`
}

func (AnnouncementImage) TableName() string {
	return "announcement_images"
}
