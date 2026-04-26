package models

import "gorm.io/gorm"

type TicketStatusHistory struct {
	gorm.Model

	TicketID    uint   `gorm:"not null;index"`
	FromStatus  string `gorm:"size:32"`
	ToStatus    string `gorm:"size:32;not null"`
	ChangedByID uint   `gorm:"not null;index"`
	ChangedBy   User   `gorm:"foreignKey:ChangedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Note        string `gorm:"type:text"`
}

func (TicketStatusHistory) TableName() string {
	return "ticket_status_history"
}
