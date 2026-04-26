package models

import "gorm.io/gorm"

type TicketAttachment struct {
	gorm.Model

	TicketID  uint  `gorm:"not null;index"`
	CommentID *uint `gorm:"index"`

	UploadedByID uint `gorm:"not null;index"`
	UploadedBy   User `gorm:"foreignKey:UploadedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	FileURL  string `gorm:"size:512;not null"`
	FileName string `gorm:"size:255;not null"`
	MimeType string `gorm:"size:128"`
	Size     int64  `gorm:"default:0"`
}

func (TicketAttachment) TableName() string {
	return "ticket_attachments"
}
