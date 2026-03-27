package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

func (uc *UseCase) CreateQuantityDiscount(ctx context.Context, dto dtos.CreateQuantityDiscountDTO) (*entities.QuantityDiscount, error) {
	if dto.MinQuantity < 1 {
		return nil, domainerrors.ErrInvalidMinQuantity
	}
	if dto.DiscountPercent <= 0 || dto.DiscountPercent > 100 {
		return nil, domainerrors.ErrInvalidDiscountPercent
	}

	discount := &entities.QuantityDiscount{
		BusinessID:      dto.BusinessID,
		ProductID:       dto.ProductID,
		MinQuantity:     dto.MinQuantity,
		DiscountPercent: dto.DiscountPercent,
		IsActive:        dto.IsActive,
		Description:     dto.Description,
	}

	return uc.repo.CreateQuantityDiscount(ctx, discount)
}
