package models

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model

	Name string `gorm:"size:120;not null;index"`
	Goal string `gorm:"type:text"`

	StartDate time.Time `gorm:"not null;index"`
	EndDate   time.Time `gorm:"not null;index"`

	Status string `gorm:"size:16;not null;index;default:'planned'"`

	CreatedByID uint `gorm:"not null;index"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Tickets []Ticket `gorm:"foreignKey:SprintID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Sprint) TableName() string {
	return "sprints"
}
