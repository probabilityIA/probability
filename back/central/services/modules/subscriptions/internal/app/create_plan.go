package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func (uc *UseCase) CreateSubscriptionType(ctx context.Context, dto dtos.CreateSubscriptionTypeDTO) (*entities.SubscriptionType, error) {
	if dto.Name == "" || dto.Code == "" || dto.Price <= 0 {
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

	subType := &entities.SubscriptionType{
		Name:          dto.Name,
		Code:          dto.Code,
		Description:   dto.Description,
		Price:         dto.Price,
		BillingPeriod: billingPeriod,
		Active:        true,
		ModuleCodes:   dto.ModuleCodes,
	}

	if err := uc.repo.CreateSubscriptionType(ctx, subType); err != nil {
		return nil, err
	}

	return subType, nil
}
