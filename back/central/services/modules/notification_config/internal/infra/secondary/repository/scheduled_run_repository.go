package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type scheduledRunRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *scheduledRunRepository) CreateRun(ctx context.Context, run *entities.ScheduledRun) error {
	model := &models.ScheduledNotificationRun{
		RuleID:     run.RuleID,
		BusinessID: run.BusinessID,
		StartedAt:  run.StartedAt,
		Status:     run.Status,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.logger.Error().Err(err).Uint("rule_id", run.RuleID).Msg("Error creating scheduled run")
		return err
	}

	run.ID = model.ID
	return nil
}

func (r *scheduledRunRepository) FinishRun(ctx context.Context, run *entities.ScheduledRun) error {
	if err := r.db.Conn(ctx).Model(&models.ScheduledNotificationRun{}).
		Where("id = ?", run.ID).
		Updates(map[string]any{
			"finished_at":      run.FinishedAt,
			"status":           run.Status,
			"matched_count":    run.MatchedCount,
			"queued_count":     run.QueuedCount,
			"skipped_count":    run.SkippedCount,
			"failed_count":     run.FailedCount,
			"cap_reached_flag": run.CapReachedFlag,
			"error_message":    run.ErrorMessage,
		}).Error; err != nil {
		r.logger.Error().Err(err).Uint("run_id", run.ID).Msg("Error finishing scheduled run")
		return err
	}
	return nil
}

func (r *scheduledRunRepository) ListRuns(ctx context.Context, ruleID uint, page, pageSize int) ([]entities.ScheduledRun, int64, error) {
	var rows []models.ScheduledNotificationRun
	var total int64

	query := r.db.Conn(ctx).Model(&models.ScheduledNotificationRun{}).Where("rule_id = ?", ruleID)

	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error counting scheduled runs")
		return nil, 0, err
	}

	if err := query.Order("started_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing scheduled runs")
		return nil, 0, err
	}

	out := make([]entities.ScheduledRun, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.ScheduledRun{
			ID:             row.ID,
			RuleID:         row.RuleID,
			BusinessID:     row.BusinessID,
			StartedAt:      row.StartedAt,
			FinishedAt:     row.FinishedAt,
			Status:         row.Status,
			MatchedCount:   row.MatchedCount,
			QueuedCount:    row.QueuedCount,
			SkippedCount:   row.SkippedCount,
			FailedCount:    row.FailedCount,
			CapReachedFlag: row.CapReachedFlag,
			ErrorMessage:   row.ErrorMessage,
		})
	}

	return out, total, nil
}
