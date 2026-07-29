package request

import (
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
)

type LeadAnswerRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LeadRecommendationRequest struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  string `json:"type"`
}

type CreateLeadRequest struct {
	Name            string                      `json:"name"`
	Email           string                      `json:"email"`
	Phone           string                      `json:"phone"`
	SurveySlug      string                      `json:"survey_slug"`
	ScoreTotal      int                         `json:"score_total"`
	ScoreMax        int                         `json:"score_max"`
	Level           string                      `json:"level"`
	Answers         []LeadAnswerRequest         `json:"answers"`
	Recommendations []LeadRecommendationRequest `json:"recommendations"`
}

func (r *CreateLeadRequest) ToDTO() dtos.CreateLeadDTO {
	answers := make([]entities.LeadAnswer, 0, len(r.Answers))
	for _, a := range r.Answers {
		answers = append(answers, entities.LeadAnswer{Question: a.Question, Answer: a.Answer})
	}

	recs := make([]entities.LeadRecommendation, 0, len(r.Recommendations))
	for _, rec := range r.Recommendations {
		recs = append(recs, entities.LeadRecommendation{Title: rec.Title, Desc: rec.Desc, Type: rec.Type})
	}

	return dtos.CreateLeadDTO{
		Name:            r.Name,
		Email:           r.Email,
		Phone:           r.Phone,
		SurveySlug:      r.SurveySlug,
		ScoreTotal:      r.ScoreTotal,
		ScoreMax:        r.ScoreMax,
		Level:           r.Level,
		Answers:         answers,
		Recommendations: recs,
	}
}
