package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
)

func (uc *UseCase) Create(ctx context.Context, dto dtos.CreateShippingMarginDTO) (*entities.ShippingMargin, error) {
	code := strings.ToLower(strings.TrimSpace(dto.CarrierCode))
	if code == "" {
		return nil, domainerrors.ErrInvalidCarrierCode
	}
	if dto.MarginAmount < 0 || dto.InsuranceMargin < 0 {
		return nil, domainerrors.ErrInvalidMargin
	}

	exists, err := uc.repo.ExistsByCarrier(ctx, dto.BusinessID, code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCarrier
	}

	m := &entities.ShippingMargin{
		BusinessID:      dto.BusinessID,
		CarrierCode:     code,
		CarrierName:     dto.CarrierName,
		MarginAmount:    dto.MarginAmount,
		InsuranceMargin: dto.InsuranceMargin,
		IsActive:        dto.IsActive,
	}
	created, err := uc.repo.Create(ctx, m)
	if err != nil {
		return nil, err
	}
	if uc.cache != nil {
		_ = uc.cache.Upsert(ctx, created)
	}
	return created, nil
}
