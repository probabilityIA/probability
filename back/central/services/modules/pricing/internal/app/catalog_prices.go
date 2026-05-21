package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
)

func validateTarget(target dtos.CatalogPriceTarget) error {
	hasGroup := target.ClientGroupID != nil
	hasClient := target.ClientID != nil
	if !hasGroup && !hasClient {
		return domainerrors.ErrTargetRequired
	}
	if hasGroup && hasClient {
		return domainerrors.ErrTargetAmbiguous
	}
	return nil
}

func (uc *UseCase) ListCatalogPrices(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error) {
	if err := validateTarget(params.Target); err != nil {
		return nil, 0, err
	}
	return uc.repo.ListCatalogPrices(ctx, params)
}

func (uc *UseCase) SaveCatalogPrices(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
	if err := validateTarget(dto.Target); err != nil {
		return err
	}
	for _, item := range dto.Items {
		if item.Price != nil && *item.Price < 0 {
			return domainerrors.ErrInvalidPrice
		}
	}
	return uc.repo.SaveCatalogPrices(ctx, dto)
}
