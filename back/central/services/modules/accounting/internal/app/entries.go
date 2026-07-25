package app

import (
	"context"
	"math"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (uc *UseCase) ListEntries(ctx context.Context, params dtos.ListEntriesParams) ([]entities.Entry, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	return uc.repo.ListEntries(ctx, params)
}

func (uc *UseCase) CreateManualEntry(ctx context.Context, dto dtos.CreateEntryDTO) (*entities.Entry, error) {
	if dto.Amount <= 0 {
		return nil, domainerrors.ErrInvalidAmount
	}
	concept, err := uc.repo.GetConceptByID(ctx, dto.ConceptID)
	if err != nil {
		return nil, err
	}
	if !concept.IsActive {
		return nil, domainerrors.ErrConceptInactive
	}
	entry := &entities.Entry{
		ConceptID:   concept.ID,
		BusinessID:  dto.BusinessID,
		EntryDate:   dto.EntryDate,
		Amount:      dto.Amount,
		Kind:        concept.Kind,
		SourceType:  entities.SourceManual,
		SourceID:    "",
		Description: strings.TrimSpace(dto.Description),
		IsAutomatic: false,
	}
	applyTaxes(entry, concept)
	if _, err := uc.repo.CreateEntry(ctx, entry); err != nil {
		return nil, err
	}
	return uc.repo.GetEntryByID(ctx, entry.ID)
}

func (uc *UseCase) DeleteManualEntry(ctx context.Context, id uint) error {
	entry, err := uc.repo.GetEntryByID(ctx, id)
	if err != nil {
		return err
	}
	if entry.IsAutomatic || entry.SourceType != entities.SourceManual {
		return domainerrors.ErrEntryNotManual
	}
	return uc.repo.DeleteEntry(ctx, id)
}

func applyTaxes(entry *entities.Entry, concept *entities.Concept) {
	entry.TaxDetail = []entities.TaxLine{}
	entry.TaxTotal = 0
	for _, ct := range concept.Taxes {
		if !ct.IsActive {
			continue
		}
		amount := round2(entry.Amount * ct.RatePercent / 100)
		entry.TaxDetail = append(entry.TaxDetail, entities.TaxLine{
			TaxCode:     ct.TaxCode,
			TaxName:     ct.TaxName,
			RatePercent: ct.RatePercent,
			Amount:      amount,
		})
		entry.TaxTotal += amount
	}
	entry.TaxTotal = round2(entry.TaxTotal)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
