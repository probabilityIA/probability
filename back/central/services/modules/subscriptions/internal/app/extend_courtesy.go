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

	meta, err := uc.repo.GetBusinessSubscriptionMeta(ctx, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, errs.ErrSubscriptionNotFound
	}

	now := time.Now()
	currentCutoff := effectiveCutoffDate(meta, now)

	base := currentCutoff
	if base.Before(now) {
		base = now
	}
	newCourtesyUntil := base.AddDate(0, 0, dto.Days)

	if err := uc.repo.UpdateBusinessSubscriptionCourtesyUntil(ctx, dto.BusinessID, newCourtesyUntil); err != nil {
		return nil, err
	}

	uc.recordAudit(ctx, dto.BusinessID, actorUserID, entities.AuditActionCourtesyExtended,
		fmt.Sprintf("extendio %d dias de cortesia: bloqueo pospuesto hasta %s. Motivo: %s", dto.Days, newCourtesyUntil.Format("2006-01-02"), dto.Reason))

	if err := uc.reactivateIfUnblocked(ctx, dto.BusinessID, actorUserID, meta.Status, newCourtesyUntil, now); err != nil {
		return nil, err
	}

	return latest, nil
}

func effectiveCutoffDate(meta *entities.BusinessSubscriptionMeta, now time.Time) time.Time {
	cutoff := now
	if meta.EndDate != nil {
		cutoff = *meta.EndDate
		if meta.CutoffDay != nil {
			cutoff = nextCutoffOnOrAfter(cutoff, *meta.CutoffDay)
		}
	}
	if meta.CourtesyUntil != nil && meta.CourtesyUntil.After(cutoff) {
		cutoff = *meta.CourtesyUntil
	}
	return cutoff
}

func (uc *UseCase) reactivateIfUnblocked(ctx context.Context, businessID, actorUserID uint, currentStatus string, unblockedUntil, now time.Time) error {
	if currentStatus != entities.BusinessStatusExpired {
		return nil
	}
	if !unblockedUntil.After(now) {
		return nil
	}
	if err := uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, entities.BusinessStatusActive, nil); err != nil {
		return err
	}
	uc.recordAudit(ctx, businessID, actorUserID, entities.AuditActionSubscriptionReactivated,
		fmt.Sprintf("cuenta reactivada automaticamente: nueva fecha de bloqueo %s", unblockedUntil.Format("2006-01-02")))
	return nil
}
