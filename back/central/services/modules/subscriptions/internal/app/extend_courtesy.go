package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
)

func (uc *UseCase) ExtendCourtesy(ctx context.Context, dto dtos.ExtendCourtesyDTO, actorUserID uint) (*entities.BusinessSubscription, error) {
	if dto.Days <= 0 {
		return nil, errs.ErrInvalidDays
	}

	latest, err := uc.repo.GetLatestByBusinessID(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if latest == nil {
		return nil, errs.ErrSubscriptionNotFound
	}

	base := latest.EndDate
	if base.Before(time.Now()) {
		base = time.Now()
	}
	newEnd := base.AddDate(0, 0, dto.Days)

	if err := uc.repo.UpdateSubscriptionDates(ctx, latest.ID, latest.StartDate, newEnd); err != nil {
		return nil, err
	}
	if err := uc.repo.UpdateBusinessSubscriptionEndDate(ctx, dto.BusinessID, newEnd); err != nil {
		return nil, err
	}

	latest.EndDate = newEnd
	uc.recordAudit(ctx, dto.BusinessID, actorUserID, entities.AuditActionCourtesyExtended,
		fmt.Sprintf("extendio %d dias de cortesia hasta %s. Motivo: %s", dto.Days, newEnd.Format("2006-01-02"), dto.Reason))

	return latest, nil
}
