package models

import "gorm.io/gorm"

// ContactSubmission stores contact form submissions from public business pages
type ContactSubmission struct {
	gorm.Model
	BusinessID uint   `gorm:"not null;index"`
	Name       string `gorm:"size:255;not null"`
	Email      string `gorm:"size:255"`
	Phone      string `gorm:"size:50"`
	Message    string `gorm:"type:text;not null"`
	Source     string `gorm:"size:50;default:'website'"` // website, whatsapp
	Status     string `gorm:"size:20;default:'new'"`     // new, read, replied
}
