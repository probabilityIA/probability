package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	gorm.Model
	BusinessID    *uint                `gorm:"index"`
	CategoryID    uint                 `gorm:"not null;index"`
	Category      AnnouncementCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Title         string               `gorm:"size:255;not null"`
	Message       string               `gorm:"type:text"`
	DisplayType   string               `gorm:"size:20;not null;index"`
	FrequencyType string               `gorm:"size:30;not null"`
	Priority      int                  `gorm:"default:0;index"`
	IsGlobal      bool                 `gorm:"default:false;index"`
	Status        string               `gorm:"size:20;not null;default:'draft';index"`
	StartsAt      *time.Time           `gorm:"index"`
	EndsAt        *time.Time           `gorm:"index"`
	ForceRedisplay bool               `gorm:"default:false"`
	CreatedByID   uint                 `gorm:"not null;index"`
	CreatedBy     User                 `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Images  []AnnouncementImage  `gorm:"foreignKey:AnnouncementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Links   []AnnouncementLink   `gorm:"foreignKey:AnnouncementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Targets []AnnouncementTarget `gorm:"foreignKey:AnnouncementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Announcement) TableName() string {
	return "announcements"
}
