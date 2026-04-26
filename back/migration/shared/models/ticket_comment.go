package models

import "gorm.io/gorm"

type TicketComment struct {
	gorm.Model

	TicketID uint   `gorm:"not null;index"`
	UserID   uint   `gorm:"not null;index"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Body     string `gorm:"type:text;not null"`
	IsInternal bool `gorm:"default:false;index"`

	Attachments []TicketAttachment `gorm:"foreignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (TicketComment) TableName() string {
	return "ticket_comments"
}
