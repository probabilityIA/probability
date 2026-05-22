package models

import (
	"time"

	"gorm.io/datatypes"
)

type Geozone struct {
	ID         uint           `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time     `gorm:"index"`
	BusinessID uint           `gorm:"not null;index;default:0"`
	ParentID   *uint          `gorm:"index"`
	Type       string         `gorm:"size:32;not null;index"`
	Code       *string        `gorm:"size:64;index"`
	Name       string         `gorm:"size:255;not null"`
	Geometry   datatypes.JSON `gorm:"type:jsonb"`
	Centroid   datatypes.JSON `gorm:"type:jsonb"`
	Properties datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	IsActive   bool           `gorm:"not null;default:true"`

	Parent   *Geozone  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Children []Geozone `gorm:"foreignKey:ParentID"`
}

func (Geozone) TableName() string {
	return "geozones"
}
