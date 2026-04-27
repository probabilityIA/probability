package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model

	Code string `gorm:"size:20;uniqueIndex;not null"`

	BusinessID *uint     `gorm:"index"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	CreatedByID uint `gorm:"not null;index"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	AssignedToID *uint `gorm:"index"`
	AssignedTo   *User `gorm:"foreignKey:AssignedToID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text;not null"`

	Type     string `gorm:"size:32;not null;index;default:'support'"`
	Category string `gorm:"size:64;index"`
	Priority string `gorm:"size:16;not null;index;default:'medium'"`
	Status   string `gorm:"size:32;not null;index;default:'open'"`
	Source   string `gorm:"size:16;not null;index;default:'internal'"`
	Severity string `gorm:"size:16;index"`
	Area     string `gorm:"size:32;index;default:'soporte'"`

	EscalatedToDev bool       `gorm:"default:false;index"`
	EscalatedAt    *time.Time `gorm:"index"`

	DueDate    *time.Time `gorm:"index"`
	ResolvedAt *time.Time `gorm:"index"`
	ClosedAt   *time.Time `gorm:"index"`

	Comments    []TicketComment    `gorm:"foreignKey:TicketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Attachments []TicketAttachment `gorm:"foreignKey:TicketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	History     []TicketStatusHistory `gorm:"foreignKey:TicketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Ticket) TableName() string {
	return "tickets"
}
