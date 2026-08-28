package models

import "time"

type LegalAcceptance struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null;index;uniqueIndex:idx_legal_user_document,priority:1"`
	LegalDocumentID uint      `gorm:"not null;index;uniqueIndex:idx_legal_user_document,priority:2"`
	BusinessID      *uint     `gorm:"index"`
	DocumentCode    string    `gorm:"size:50;not null"`
	DocumentVersion string    `gorm:"size:20;not null"`
	DocumentSHA256  string    `gorm:"size:64;not null"`
	AcceptedAt      time.Time `gorm:"not null;index"`
	IPAddress       string    `gorm:"size:64"`
	UserAgent       string    `gorm:"size:500"`
	Method          string    `gorm:"size:30;not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`

	User          User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	LegalDocument LegalDocument `gorm:"foreignKey:LegalDocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (LegalAcceptance) TableName() string {
	return "legal_acceptances"
}
