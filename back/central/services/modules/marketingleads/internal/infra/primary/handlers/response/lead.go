package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
)

type LeadResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	SurveySlug        string    `json:"survey_slug"`
	ScoreTotal        int       `json:"score_total"`
	ScoreMax          int       `json:"score_max"`
	Level             string    `json:"level"`
	WhatsAppMessageID string    `json:"whats_app_message_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

func FromEntity(l *entities.MarketingLead) LeadResponse {
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
