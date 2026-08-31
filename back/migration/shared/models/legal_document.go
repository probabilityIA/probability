package models

import (
	"time"

	"gorm.io/gorm"
)

type LegalDocument struct {
	gorm.Model
	Code          string    `gorm:"size:50;not null;uniqueIndex:idx_legal_code_version,priority:1"`
	Version       string    `gorm:"size:20;not null;uniqueIndex:idx_legal_code_version,priority:2"`
	Title         string    `gorm:"size:255;not null"`
	FileURL       string    `gorm:"size:500;not null"`
	SHA256        string    `gorm:"size:64;not null"`
	ContentHTML   string    `gorm:"type:text"`
	EffectiveFrom time.Time `gorm:"not null"`
	IsActive      bool      `gorm:"not null;default:false;index"`
	RequiresStaff bool      `gorm:"not null;default:true"`
}

func (LegalDocument) TableName() string {
	return "legal_documents"
}
