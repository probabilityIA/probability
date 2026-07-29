package request

import "github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/dtos"

type CreateLeadRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	SurveySlug string `json:"survey_slug"`
	ScoreTotal int    `json:"score_total"`
	ScoreMax   int    `json:"score_max"`
	Level      string `json:"level"`
}

func (r *CreateLeadRequest) ToDTO() dtos.CreateLeadDTO {
	return dtos.CreateLeadDTO{
		Name:       r.Name,
		Email:      r.Email,
		Phone:      r.Phone,
		SurveySlug: r.SurveySlug,
		ScoreTotal: r.ScoreTotal,
		ScoreMax:   r.ScoreMax,
		Level:      r.Level,
	}
}
