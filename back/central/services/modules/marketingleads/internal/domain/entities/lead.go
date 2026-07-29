package entities

import "time"

type MarketingLead struct {
	ID                uint
	Name              string
	Email             string
	Phone             string
	SurveySlug        string
	ScoreTotal        int
	ScoreMax          int
	Level             string
	WhatsAppMessageID string
	CreatedAt         time.Time
}
