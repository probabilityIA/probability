package models

import "gorm.io/gorm"

type AnnouncementImage struct {
	gorm.Model
	AnnouncementID uint   `gorm:"not null;index"`
	ImageURL       string `gorm:"size:500;not null"`
	SortOrder      int    `gorm:"default:0"`
}

func (AnnouncementImage) TableName() string {
	return "announcement_images"
}
