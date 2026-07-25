package request

type CreateConceptRequest struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	Kind         string `json:"kind" binding:"required"`
	IsRealIncome *bool  `json:"is_real_income"`
}

type UpdateConceptRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	IsRealIncome *bool  `json:"is_real_income"`
	IsActive     *bool  `json:"is_active"`
}

type CreateTaxRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	RatePercent float64 `json:"rate_percent"`
	Kind        string  `json:"kind"`
}

type UpdateTaxRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	RatePercent float64 `json:"rate_percent"`
	Kind        string  `json:"kind"`
	IsActive    *bool   `json:"is_active"`
}

type SetConceptTaxRequest struct {
	IsActive bool `json:"is_active"`
}

type CreateEntryRequest struct {
	ConceptID   uint    `json:"concept_id" binding:"required"`
	BusinessID  *uint   `json:"related_business_id"`
	EntryDate   string  `json:"entry_date" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}
