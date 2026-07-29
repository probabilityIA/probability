package dtos

type CreateLeadDTO struct {
	Name       string
	Email      string
	Phone      string
	SurveySlug string
	ScoreTotal int
	ScoreMax   int
	Level      string
}
