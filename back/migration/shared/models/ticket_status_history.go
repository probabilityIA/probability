package models

import "gorm.io/gorm"

type TicketStatusHistory struct {
	gorm.Model

	TicketID    uint   `gorm:"not null;index"`
	ChangeType  string `gorm:"size:16;not null;default:'status';index"`
	FromStatus  string `gorm:"size:32"`
	ToStatus    string `gorm:"size:32"`
	FromArea    string `gorm:"size:32"`
	ToArea      string `gorm:"size:32"`
	ChangedByID uint   `gorm:"not null;index"`
	ChangedBy   User   `gorm:"foreignKey:ChangedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Note        string `gorm:"type:text"`
}

func (TicketStatusHistory) TableName() string {
	return "ticket_status_history"
}
