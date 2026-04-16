package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

func (uc *UseCase) UpdateClientPricingRule(ctx context.Context, dto dtos.UpdateClientPricingRuleDTO) (*entities.ClientPricingRule, error) {
	if dto.AdjustmentType != "percentage" && dto.AdjustmentType != "fixed" {
		return nil, domainerrors.ErrInvalidAdjustmentType
	}

	rule := &entities.ClientPricingRule{
		ID:              dto.ID,
		BusinessID:      dto.BusinessID,
		ClientID:        dto.ClientID,
		ProductID:       dto.ProductID,
		AdjustmentType:  dto.AdjustmentType,
		AdjustmentValue: dto.AdjustmentValue,
		IsActive:        dto.IsActive,
		Priority:        dto.Priority,
		Description:     dto.Description,
	}

	return uc.repo.UpdateClientPricingRule(ctx, rule)
}
