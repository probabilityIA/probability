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

const detailInsertBatch = 500

func (r *repository) Save(ctx context.Context, run *domain.SyncRun) error {
	detail := datatypes.JSON(nil)

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
		ChannelNoSKU:      run.ChannelNoSKU,
		SKUChanged:        run.SKUChanged,
		SKUTypo:           run.SKUTypo,
		SKUSpacing:        run.SKUSpacing,
		Status:            run.Status,
		Message:           run.Message,
		Detail:            detail,
	}

	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := tx.Create(&row).Error; createErr != nil {
				return createErr
			}
			run.ID = row.ID
			return replaceDetailItems(tx, row.ID, run.Detail)
		}

		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
		if updateErr := tx.Model(&models.IntegrationSyncRun{}).
			Where("id = ?", existing.ID).
			Select("correlation_id", "started_at", "finished_at", "total", "updated", "unchanged",
				"skipped", "failed", "matched", "not_associated", "only_in_probability",
				"only_in_channel", "channel_no_sku", "sku_changed", "sku_typo", "sku_spacing",
				"status", "message", "detail", "updated_at").
			Updates(&row).Error; updateErr != nil {
			return updateErr
		}
		run.ID = existing.ID
		return replaceDetailItems(tx, existing.ID, run.Detail)
	})
}

func replaceDetailItems(tx *gorm.DB, runID uint, detail []domain.DetailItem) error {
	if err := tx.Where("run_id = ?", runID).Delete(&models.IntegrationSyncRunItem{}).Error; err != nil {
		return err
	}
	if len(detail) == 0 {
		return nil
	}

	rows := make([]models.IntegrationSyncRunItem, 0, len(detail))
	for _, item := range detail {
		rows = append(rows, models.IntegrationSyncRunItem{
			RunID:        runID,
			GroupCode:    truncate(item.Group, 32),
			SKU:          truncate(item.SKU, 120),
			Label:        truncate(item.Label, 300),
			Tone:         truncate(item.Tone, 16),
			ParentRef:    truncate(item.ParentRef, 64),
			ParentLabel:  truncate(item.ParentLabel, 300),
			VariantLabel: truncate(item.VariantLabel, 200),

			CounterpartSKU:  truncate(item.CounterpartSKU, 120),
			CounterpartName: truncate(item.CounterpartName, 300),
			ChannelQty:      item.ChannelQty,
			OwnQty:          item.OwnQty,
			FixSide:         truncate(item.FixSide, 16),
			Pattern:         truncate(item.Pattern, 24),
		})
	}
	return tx.CreateInBatches(&rows, detailInsertBatch).Error
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func (r *repository) ListDetail(ctx context.Context, query domain.DetailQuery) (*domain.DetailPage, error) {
	page := &domain.DetailPage{
		Items:    []domain.DetailItem{},
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	var run models.IntegrationSyncRun
	err := r.db.Conn(ctx).
		Select("id").
		Where("business_id = ? AND integration_id = ? AND kind = ?", query.BusinessID, query.IntegrationID, query.Kind).
		First(&run).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return page, nil
	}
	if err != nil {
		return nil, err
	}

	base := r.db.Conn(ctx).
		Model(&models.IntegrationSyncRunItem{}).
		Where("run_id = ?", run.ID)

	if query.Group != "" {
		base = base.Where("group_code = ?", query.Group)
	}
	if query.Search != "" {
		like := "%" + query.Search + "%"
		base = base.Where("sku ILIKE ? OR label ILIKE ?", like, like)
	}

	if err := base.Count(&page.Total).Error; err != nil {
		return nil, err
	}

	var rows []models.IntegrationSyncRunItem
	if err := base.
		Order("id ASC").
		Offset(query.Offset()).
		Limit(query.PageSize).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		page.Items = append(page.Items, domain.DetailItem{
			SKU:          row.SKU,
			Label:        row.Label,
			Tone:         row.Tone,
			Group:        row.GroupCode,
			ParentRef:    row.ParentRef,
			ParentLabel:  row.ParentLabel,
			VariantLabel: row.VariantLabel,
		})
	}

	page.TotalPages = int((page.Total + int64(query.PageSize) - 1) / int64(query.PageSize))
	return page, nil
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
			ChannelNoSKU:      row.ChannelNoSKU,
			SKUChanged:        row.SKUChanged,
			SKUTypo:           row.SKUTypo,
			SKUSpacing:        row.SKUSpacing,
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
