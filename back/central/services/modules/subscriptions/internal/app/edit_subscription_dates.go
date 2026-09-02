package app

import (
	"context"
	"fmt"
	"time"

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

	meta, err := uc.repo.GetBusinessSubscriptionMeta(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, errs.ErrSubscriptionNotFound
	}

	if err := uc.repo.UpdateSubscriptionDates(ctx, latest.ID, dto.StartDate, dto.EndDate); err != nil {
		return nil, err
	}
	if err := uc.repo.UpdateBusinessSubscriptionEndDate(ctx, dto.BusinessID, dto.EndDate); err != nil {
		return nil, err
	}

	cutoffDay := meta.CutoffDay
	auditMsg := fmt.Sprintf("edito las fechas del ciclo: %s -> %s", dto.StartDate.Format("2006-01-02"), dto.EndDate.Format("2006-01-02"))
	if dto.CutoffDay != nil {
		if err := uc.repo.UpdateBusinessSubscriptionCutoffDay(ctx, dto.BusinessID, *dto.CutoffDay); err != nil {
			return nil, err
		}
		cutoffDay = dto.CutoffDay
		auditMsg = fmt.Sprintf("%s, dia de corte: %d", auditMsg, *dto.CutoffDay)
	}

	latest.StartDate = dto.StartDate
	latest.EndDate = dto.EndDate
	uc.recordAudit(ctx, dto.BusinessID, actorUserID, entities.AuditActionDatesEdited, auditMsg)

	now := time.Now()
	newEndDate := dto.EndDate
	newCutoff := effectiveCutoffDate(&entities.BusinessSubscriptionMeta{
		EndDate:       &newEndDate,
		CutoffDay:     cutoffDay,
		CourtesyUntil: meta.CourtesyUntil,
	}, now)
	if err := uc.reactivateIfUnblocked(ctx, dto.BusinessID, actorUserID, meta.Status, newCutoff, now); err != nil {
		return nil, err
	}

	return latest, nil
}
