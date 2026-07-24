package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (uc *UseCase) ListTaxes(ctx context.Context) ([]entities.Tax, error) {
	return uc.repo.ListTaxes(ctx)
}

func (uc *UseCase) CreateTax(ctx context.Context, dto dtos.CreateTaxDTO) (*entities.Tax, error) {
	dto.Code = strings.ToUpper(strings.TrimSpace(dto.Code))
	dto.Name = strings.TrimSpace(dto.Name)
	if dto.Code == "" {
		return nil, domainerrors.ErrCodeRequired
	}
	if dto.Name == "" {
		return nil, domainerrors.ErrNameRequired
	}
	if dto.RatePercent < 0 || dto.RatePercent > 100 {
		return nil, domainerrors.ErrInvalidRate
	}
	exists, err := uc.repo.TaxCodeExists(ctx, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCode
	}
	return uc.repo.CreateTax(ctx, dto)
}

func (uc *UseCase) UpdateTax(ctx context.Context, dto dtos.UpdateTaxDTO) (*entities.Tax, error) {
	dto.Name = strings.TrimSpace(dto.Name)
	if dto.Name == "" {
		return nil, domainerrors.ErrNameRequired
	}
	if dto.RatePercent < 0 || dto.RatePercent > 100 {
		return nil, domainerrors.ErrInvalidRate
	}
	if _, err := uc.repo.GetTaxByID(ctx, dto.ID); err != nil {
		return nil, err
	}
	return uc.repo.UpdateTax(ctx, dto)
}
