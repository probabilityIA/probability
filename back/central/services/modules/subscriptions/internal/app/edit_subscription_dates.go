package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

func (uc *UseCase) EditSubscriptionDates(ctx context.Context, dto dtos.EditSubscriptionDatesDTO, actorUserID uint) (*entities.BusinessSubscription, error) {
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
	uc.recordAudit(ctx, dto.BusinessID, actorUserID, entities.AuditActionDatesEdited,
		fmt.Sprintf("edito las fechas del ciclo: %s -> %s", dto.StartDate.Format("2006-01-02"), dto.EndDate.Format("2006-01-02")))
	return latest, nil
}
