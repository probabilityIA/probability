package response

import (
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

type ConceptTaxResponse struct {
	ID          uint    `json:"id"`
	TaxID       uint    `json:"tax_id"`
	TaxCode     string  `json:"tax_code"`
	TaxName     string  `json:"tax_name"`
	RatePercent float64 `json:"rate_percent"`
	IsActive    bool    `json:"is_active"`
}

type ConceptResponse struct {
	ID           uint                 `json:"id"`
	Code         string               `json:"code"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Kind         string               `json:"kind"`
	IsRealIncome bool                 `json:"is_real_income"`
	IsAutomatic  bool                 `json:"is_automatic"`
	SourceType   string               `json:"source_type"`
	IsActive     bool                 `json:"is_active"`
	Taxes        []ConceptTaxResponse `json:"taxes"`
}

func FromConcept(e *entities.Concept) ConceptResponse {
	taxes := make([]ConceptTaxResponse, len(e.Taxes))
	for i, t := range e.Taxes {
		taxes[i] = ConceptTaxResponse{
			ID:          t.ID,
			TaxID:       t.TaxID,
			TaxCode:     t.TaxCode,
			TaxName:     t.TaxName,
			RatePercent: t.RatePercent,
			IsActive:    t.IsActive,
		}
	}
	return ConceptResponse{
		ID:           e.ID,
		Code:         e.Code,
		Name:         e.Name,
		Description:  e.Description,
		Kind:         e.Kind,
		IsRealIncome: e.IsRealIncome,
		IsAutomatic:  e.IsAutomatic,
		SourceType:   e.SourceType,
		IsActive:     e.IsActive,
		Taxes:        taxes,
	}
}

type TaxResponse struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	RatePercent float64 `json:"rate_percent"`
	Kind        string  `json:"kind"`
	IsActive    bool    `json:"is_active"`
}

func FromTax(e *entities.Tax) TaxResponse {
	return TaxResponse{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		RatePercent: e.RatePercent,
		Kind:        e.Kind,
		IsActive:    e.IsActive,
	}
}

type TaxLineResponse struct {
	TaxCode     string  `json:"tax_code"`
	TaxName     string  `json:"tax_name"`
	RatePercent float64 `json:"rate_percent"`
	Amount      float64 `json:"amount"`
}

type EntryResponse struct {
	ID           uint              `json:"id"`
	ConceptID    uint              `json:"concept_id"`
	ConceptCode  string            `json:"concept_code"`
	ConceptName  string            `json:"concept_name"`
	BusinessID   *uint             `json:"business_id"`
	BusinessName string            `json:"business_name"`
	EntryDate    string            `json:"entry_date"`
	Amount       float64           `json:"amount"`
	Kind         string            `json:"kind"`
	SourceType   string            `json:"source_type"`
	SourceID     string            `json:"source_id"`
	Description  string            `json:"description"`
	IsAutomatic  bool              `json:"is_automatic"`
	TaxTotal     float64           `json:"tax_total"`
	TaxDetail    []TaxLineResponse `json:"tax_detail"`
}

func FromEntry(e *entities.Entry) EntryResponse {
	lines := make([]TaxLineResponse, len(e.TaxDetail))
	for i, l := range e.TaxDetail {
		lines[i] = TaxLineResponse{
			TaxCode:     l.TaxCode,
			TaxName:     l.TaxName,
			RatePercent: l.RatePercent,
			Amount:      l.Amount,
		}
	}
	return EntryResponse{
		ID:           e.ID,
		ConceptID:    e.ConceptID,
		ConceptCode:  e.ConceptCode,
		ConceptName:  e.ConceptName,
		BusinessID:   e.BusinessID,
		BusinessName: e.BusinessName,
		EntryDate:    e.EntryDate.Format("2006-01-02"),
		Amount:       e.Amount,
		Kind:         e.Kind,
		SourceType:   e.SourceType,
		SourceID:     e.SourceID,
		Description:  e.Description,
		IsAutomatic:  e.IsAutomatic,
		TaxTotal:     e.TaxTotal,
		TaxDetail:    lines,
	}
}

type EntriesListResponse struct {
	Success    bool            `json:"success"`
	Data       []EntryResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
