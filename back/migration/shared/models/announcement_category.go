package models

import "gorm.io/gorm"

type AnnouncementCategory struct {
	gorm.Model
	Code  string `gorm:"size:50;not null;uniqueIndex"`
	Name  string `gorm:"size:100;not null"`
	Icon  string `gorm:"size:100"`
	Color string `gorm:"size:20"`
}

func (AnnouncementCategory) TableName() string {
	return "announcement_categories"
}
