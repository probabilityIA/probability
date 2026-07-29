package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
)

type LeadAnswerResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LeadRecommendationResponse struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  string `json:"type,omitempty"`
}

type LeadResponse struct {
	ID                uint                         `json:"id"`
	Name              string                       `json:"name"`
	Email             string                       `json:"email"`
	Phone             string                       `json:"phone"`
	SurveySlug        string                       `json:"survey_slug"`
	ScoreTotal        int                          `json:"score_total"`
	ScoreMax          int                          `json:"score_max"`
	Level             string                       `json:"level"`
	WhatsAppMessageID string                       `json:"whats_app_message_id,omitempty"`
	Answers           []LeadAnswerResponse         `json:"answers"`
	Recommendations   []LeadRecommendationResponse `json:"recommendations"`
	CreatedAt         time.Time                    `json:"created_at"`
}

func FromEntity(l *entities.MarketingLead) LeadResponse {
	answers := make([]LeadAnswerResponse, 0, len(l.Answers))
	for _, a := range l.Answers {
		answers = append(answers, LeadAnswerResponse{Question: a.Question, Answer: a.Answer})
	}

	recs := make([]LeadRecommendationResponse, 0, len(l.Recommendations))
	for _, rec := range l.Recommendations {
		recs = append(recs, LeadRecommendationResponse{Title: rec.Title, Desc: rec.Desc, Type: rec.Type})
	}

	return LeadResponse{
		ID:                l.ID,
		Name:              l.Name,
		Email:             l.Email,
		Phone:             l.Phone,
		SurveySlug:        l.SurveySlug,
		ScoreTotal:        l.ScoreTotal,
		ScoreMax:          l.ScoreMax,
		Level:             l.Level,
		WhatsAppMessageID: l.WhatsAppMessageID,
		Answers:           answers,
		Recommendations:   recs,
		CreatedAt:         l.CreatedAt,
	}
}

func FromEntities(leads []entities.MarketingLead) []LeadResponse {
	out := make([]LeadResponse, 0, len(leads))
	for i := range leads {
		out = append(out, FromEntity(&leads[i]))
	}
	return out
}
