package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
)

func (uc *UseCase) Update(ctx context.Context, dto dtos.UpdateShippingMarginDTO) (*entities.ShippingMargin, error) {
	if dto.MarginAmount < 0 || dto.InsuranceMargin < 0 {
		return nil, domainerrors.ErrInvalidMargin
	}

	existing, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	existing.CarrierName = dto.CarrierName
	existing.MarginAmount = dto.MarginAmount
	existing.InsuranceMargin = dto.InsuranceMargin
	existing.IsActive = dto.IsActive

	updated, err := uc.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	if uc.cache != nil {
		if updated.IsActive {
			_ = uc.cache.Upsert(ctx, updated)
		} else {
			_ = uc.cache.Delete(ctx, updated.BusinessID, updated.CarrierCode)
		}
	}

	return updated, nil
}
