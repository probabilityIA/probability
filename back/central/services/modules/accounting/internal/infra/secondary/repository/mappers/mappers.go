package mappers

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type taxLineJSON struct {
	TaxCode     string  `json:"tax_code"`
	TaxName     string  `json:"tax_name"`
	RatePercent float64 `json:"rate_percent"`
	Amount      float64 `json:"amount"`
}

func ConceptToEntity(m *models.AccountingConcept, taxNames map[uint]models.AccountingTax) *entities.Concept {
	concept := &entities.Concept{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		Description:  m.Description,
		Kind:         m.Kind,
		IsRealIncome: m.IsRealIncome,
		IsAutomatic:  m.IsAutomatic,
		SourceType:   m.SourceType,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		Taxes:        []entities.ConceptTax{},
	}
	for _, ct := range m.Taxes {
		tax, ok := taxNames[ct.TaxID]
		if !ok {
			continue
		}
		concept.Taxes = append(concept.Taxes, entities.ConceptTax{
			ID:          ct.ID,
			ConceptID:   ct.ConceptID,
			TaxID:       ct.TaxID,
			TaxCode:     tax.Code,
			TaxName:     tax.Name,
			RatePercent: tax.RatePercent,
			IsActive:    ct.IsActive && tax.IsActive,
		})
	}
	return concept
}

func TaxToEntity(m *models.AccountingTax) *entities.Tax {
	return &entities.Tax{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		RatePercent: m.RatePercent,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func EntryToEntity(m *models.AccountingEntry, conceptCode, conceptName, businessName string) *entities.Entry {
	entry := &entities.Entry{
		ID:           m.ID,
		ConceptID:    m.ConceptID,
		ConceptCode:  conceptCode,
		ConceptName:  conceptName,
		BusinessID:   m.BusinessID,
		BusinessName: businessName,
		EntryDate:    m.EntryDate,
		Amount:       m.Amount,
		Kind:         m.Kind,
		SourceType:   m.SourceType,
		SourceID:     m.SourceID,
		Description:  m.Description,
		IsAutomatic:  m.IsAutomatic,
		TaxTotal:     m.TaxTotal,
		TaxDetail:    []entities.TaxLine{},
		CreatedAt:    m.CreatedAt,
	}
	if len(m.TaxDetail) > 0 {
		var lines []taxLineJSON
		if err := json.Unmarshal(m.TaxDetail, &lines); err == nil {
			for _, l := range lines {
				entry.TaxDetail = append(entry.TaxDetail, entities.TaxLine{
					TaxCode:     l.TaxCode,
					TaxName:     l.TaxName,
					RatePercent: l.RatePercent,
					Amount:      l.Amount,
				})
			}
		}
	}
	return entry
}

func TaxLinesToJSON(lines []entities.TaxLine) []byte {
	out := make([]taxLineJSON, 0, len(lines))
	for _, l := range lines {
		out = append(out, taxLineJSON{
			TaxCode:     l.TaxCode,
			TaxName:     l.TaxName,
			RatePercent: l.RatePercent,
			Amount:      l.Amount,
		})
	}
	data, err := json.Marshal(out)
	if err != nil {
		return []byte("[]")
	}
	return data
}
