package dtos

import "github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"

type CreateLeadDTO struct {
	Name            string
	Email           string
	Phone           string
	SurveySlug      string
	ScoreTotal      int
	ScoreMax        int
	Level           string
	Answers         []entities.LeadAnswer
	Recommendations []entities.LeadRecommendation
}
