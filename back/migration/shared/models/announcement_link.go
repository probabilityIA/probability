package models

import "gorm.io/gorm"

type AnnouncementLink struct {
	gorm.Model
	AnnouncementID uint   `gorm:"not null;index"`
	Label          string `gorm:"size:255;not null"`
	URL            string `gorm:"size:500;not null"`
	SortOrder      int    `gorm:"default:0"`
}

func (AnnouncementLink) TableName() string {
	return "announcement_links"
}
