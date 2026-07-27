package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

// EditSubscriptionDates corrige las fechas de la suscripcion vigente de un negocio,
// sin crear un nuevo registro de pago (a diferencia de RegisterPayment/PurchaseSubscription).
func (uc *UseCase) EditSubscriptionDates(ctx context.Context, dto dtos.EditSubscriptionDatesDTO) (*entities.BusinessSubscription, error) {
	if !dto.EndDate.After(dto.StartDate) {
		return nil, errs.ErrInvalidDateRange
	}

	latest, err := uc.repo.GetLatestByBusinessID(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if latest == nil {
		return nil, errs.ErrSubscriptionNotFound
	}

	if err := uc.repo.UpdateSubscriptionDates(ctx, latest.ID, dto.StartDate, dto.EndDate); err != nil {
		return nil, err
	}
	if err := uc.repo.UpdateBusinessSubscriptionEndDate(ctx, dto.BusinessID, dto.EndDate); err != nil {
		return nil, err
	}

	latest.StartDate = dto.StartDate
	latest.EndDate = dto.EndDate
	return latest, nil
}
