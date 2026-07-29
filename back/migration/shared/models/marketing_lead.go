package models

import "gorm.io/gorm"

// MarketingLead stores leads captured from the public diagnostic quizzes
// on the marketing site (front/website), before showing them their result.
type MarketingLead struct {
	gorm.Model
	Name              string `gorm:"size:255;not null"`
	Email             string `gorm:"size:255;not null"`
	Phone             string `gorm:"size:50;not null"`
	SurveySlug        string `gorm:"size:100;not null;index"`
	ScoreTotal        int
	ScoreMax          int
	Level             string `gorm:"size:50"`
	WhatsAppMessageID string `gorm:"size:100"`
}

func (MarketingLead) TableName() string {
	return "marketing_leads"
}
