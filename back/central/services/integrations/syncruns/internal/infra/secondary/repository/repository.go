package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *repository) Save(ctx context.Context, run *domain.SyncRun) error {
	detail := datatypes.JSON(nil)
	if len(run.Detail) > 0 {
		raw, err := json.Marshal(run.Detail)
		if err == nil {
			detail = datatypes.JSON(raw)
		}
	}

	var existing models.IntegrationSyncRun
	err := r.db.Conn(ctx).
		Where("business_id = ? AND integration_id = ? AND kind = ?", run.BusinessID, run.IntegrationID, run.Kind).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	row := models.IntegrationSyncRun{
		BusinessID:        run.BusinessID,
		IntegrationID:     run.IntegrationID,
		Kind:              run.Kind,
		CorrelationID:     run.CorrelationID,
		StartedAt:         run.StartedAt,
		FinishedAt:        run.FinishedAt,
		Total:             run.Total,
		Updated:           run.Updated,
		Unchanged:         run.Unchanged,
		Skipped:           run.Skipped,
		Failed:            run.Failed,
		Matched:           run.Matched,
		NotAssociated:     run.NotAssociated,
		OnlyInProbability: run.OnlyInProbability,
		OnlyInChannel:     run.OnlyInChannel,
		Status:            run.Status,
		Message:           run.Message,
		Detail:            detail,
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if createErr := r.db.Conn(ctx).Create(&row).Error; createErr != nil {
			return createErr
		}
		run.ID = row.ID
		return nil
	}

	row.ID = existing.ID
	row.CreatedAt = existing.CreatedAt
	if updateErr := r.db.Conn(ctx).Model(&models.IntegrationSyncRun{}).
		Where("id = ?", existing.ID).
		Select("correlation_id", "started_at", "finished_at", "total", "updated", "unchanged",
			"skipped", "failed", "matched", "not_associated", "only_in_probability",
			"only_in_channel", "status", "message", "detail", "updated_at").
		Updates(&row).Error; updateErr != nil {
		return updateErr
	}
	run.ID = existing.ID
	return nil
}

func (r *repository) ListLastByBusiness(ctx context.Context, businessID uint) ([]domain.SyncRun, error) {
	var rows []models.IntegrationSyncRun
	if err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Order("finished_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	runs := make([]domain.SyncRun, 0, len(rows))
	for _, row := range rows {
		run := domain.SyncRun{
			ID:                row.ID,
			BusinessID:        row.BusinessID,
			IntegrationID:     row.IntegrationID,
			Kind:              row.Kind,
			CorrelationID:     row.CorrelationID,
			StartedAt:         row.StartedAt,
			FinishedAt:        row.FinishedAt,
			Total:             row.Total,
			Updated:           row.Updated,
			Unchanged:         row.Unchanged,
			Skipped:           row.Skipped,
			Failed:            row.Failed,
			Matched:           row.Matched,
			NotAssociated:     row.NotAssociated,
			OnlyInProbability: row.OnlyInProbability,
			OnlyInChannel:     row.OnlyInChannel,
			Status:            row.Status,
			Message:           row.Message,
		}
		if len(row.Detail) > 0 {
			var detail []domain.DetailItem
			if err := json.Unmarshal(row.Detail, &detail); err == nil {
				run.Detail = detail
			}
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func (r *repository) IntegrationBelongsToBusiness(ctx context.Context, integrationID, businessID uint) (bool, error) {
	var count int64
	if err := r.db.Conn(ctx).
		Model(&models.Integration{}).
		Where("id = ? AND business_id = ?", integrationID, businessID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
