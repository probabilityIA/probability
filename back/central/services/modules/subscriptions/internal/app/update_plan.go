package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func (uc *UseCase) UpdateSubscriptionType(ctx context.Context, dto dtos.UpdateSubscriptionTypeDTO) (*entities.SubscriptionType, error) {
	existing, err := uc.repo.GetSubscriptionType(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errs.ErrSubscriptionTypeNotFound
	}

	if dto.Name == "" || dto.Price <= 0 {
		return nil, errs.ErrInvalidSubscriptionType
	}

	for _, code := range dto.ModuleCodes {
		if !moduleregistry.IsValid(code) {
			return nil, errs.ErrInvalidModuleCode
		}
	}

	billingPeriod := dto.BillingPeriod
	if billingPeriod == "" {
		billingPeriod = "monthly"
	}

	existing.Name = dto.Name
	existing.Description = dto.Description
	existing.Price = dto.Price
	existing.BillingPeriod = billingPeriod
	existing.Active = dto.Active
	existing.ModuleCodes = dto.ModuleCodes
	existing.MaxEcommerceChannels = dto.MaxEcommerceChannels

	if err := uc.repo.UpdateSubscriptionType(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}
