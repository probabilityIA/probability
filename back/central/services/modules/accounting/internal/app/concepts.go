package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (uc *UseCase) ListConcepts(ctx context.Context) ([]entities.Concept, error) {
	return uc.repo.ListConcepts(ctx)
}

func (uc *UseCase) CreateConcept(ctx context.Context, dto dtos.CreateConceptDTO) (*entities.Concept, error) {
	dto.Code = strings.ToUpper(strings.TrimSpace(dto.Code))
	dto.Name = strings.TrimSpace(dto.Name)
	if dto.Code == "" {
		return nil, domainerrors.ErrCodeRequired
	}
	if dto.Name == "" {
		return nil, domainerrors.ErrNameRequired
	}
	if dto.Kind != entities.KindIncome && dto.Kind != entities.KindExpense {
		return nil, domainerrors.ErrInvalidKind
	}
	exists, err := uc.repo.ConceptCodeExists(ctx, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCode
	}
	return uc.repo.CreateConcept(ctx, dto)
}

func (uc *UseCase) UpdateConcept(ctx context.Context, dto dtos.UpdateConceptDTO) (*entities.Concept, error) {
	dto.Name = strings.TrimSpace(dto.Name)
	if dto.Name == "" {
		return nil, domainerrors.ErrNameRequired
	}
	if _, err := uc.repo.GetConceptByID(ctx, dto.ID); err != nil {
		return nil, err
	}
	return uc.repo.UpdateConcept(ctx, dto)
}

func (uc *UseCase) SetConceptTax(ctx context.Context, dto dtos.SetConceptTaxDTO) error {
	if _, err := uc.repo.GetConceptByID(ctx, dto.ConceptID); err != nil {
		return err
	}
	tax, err := uc.repo.GetTaxByID(ctx, dto.TaxID)
	if err != nil {
		return err
	}
	if dto.IsActive && !tax.IsActive {
		return domainerrors.ErrConceptTaxNotAllowed
	}
	return uc.repo.SetConceptTax(ctx, dto)
}
