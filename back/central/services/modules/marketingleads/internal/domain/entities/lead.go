package entities

import "time"

type LeadAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LeadRecommendation struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  string `json:"type"`
}

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
	Answers           []LeadAnswer
	Recommendations   []LeadRecommendation
	CreatedAt         time.Time
}
