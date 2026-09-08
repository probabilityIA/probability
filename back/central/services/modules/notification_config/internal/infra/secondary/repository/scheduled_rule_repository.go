package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type scheduledRuleRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *scheduledRuleRepository) CreateRule(ctx context.Context, rule *entities.ScheduledRule) error {
	model, err := ruleToModel(rule)
	if err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.logger.Error().Err(err).Str("name", rule.Name).Msg("Error creating scheduled rule")
		return err
	}

	rule.ID = model.ID
	rule.CreatedAt = model.CreatedAt
	rule.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *scheduledRuleRepository) UpdateRule(ctx context.Context, rule *entities.ScheduledRule) error {
	model, err := ruleToModel(rule)
	if err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Model(&models.ScheduledNotificationRule{}).
		Where("id = ? AND business_id = ?", rule.ID, rule.BusinessID).
		Updates(map[string]any{
			"integration_id":       model.IntegrationID,
			"notification_type_id": model.NotificationTypeID,
			"whatsapp_template_id": model.WhatsappTemplateID,
			"name":                 model.Name,
			"description":          model.Description,
			"segment_type":         model.SegmentType,
			"segment_params":       model.SegmentParams,
			"timezone":             model.Timezone,
			"frequency_minutes":    model.FrequencyMinutes,
			"send_window_start":    model.SendWindowStart,
			"send_window_end":      model.SendWindowEnd,
			"cooldown_days":        model.CooldownDays,
			"daily_send_cap":       model.DailySendCap,
			"batch_size_cap":       model.BatchSizeCap,
			"requires_opt_in":      model.RequiresOptIn,
			"enabled":              model.Enabled,
			"next_run_at":          model.NextRunAt,
		}).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", rule.ID).Msg("Error updating scheduled rule")
		return err
	}

	return nil
}

func (r *scheduledRuleRepository) GetRuleByID(ctx context.Context, id uint) (*entities.ScheduledRule, error) {
	var model models.ScheduledNotificationRule

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error getting scheduled rule")
		return nil, err
	}

	return ruleToDomain(&model)
}

func (r *scheduledRuleRepository) ListRules(ctx context.Context, businessID uint, page, pageSize int) ([]entities.ScheduledRule, int64, error) {
	var rows []models.ScheduledNotificationRule
	var total int64

	query := r.db.Conn(ctx).Model(&models.ScheduledNotificationRule{}).Where("business_id = ?", businessID)

	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error counting scheduled rules")
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing scheduled rules")
		return nil, 0, err
	}

	out := make([]entities.ScheduledRule, 0, len(rows))
	for i := range rows {
		item, err := ruleToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *item)
	}

	return out, total, nil
}

func (r *scheduledRuleRepository) DeleteRule(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.ScheduledNotificationRule{}, id).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error deleting scheduled rule")
		return err
	}
	return nil
}

func (r *scheduledRuleRepository) ListDueRules(ctx context.Context, now time.Time, limit int) ([]entities.ScheduledRule, error) {
	var rows []models.ScheduledNotificationRule

	if err := r.db.Conn(ctx).
		Where("enabled = true AND (next_run_at IS NULL OR next_run_at <= ?)", now).
		Order("next_run_at ASC NULLS FIRST").
		Limit(limit).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing due scheduled rules")
		return nil, err
	}

	out := make([]entities.ScheduledRule, 0, len(rows))
	for i := range rows {
		item, err := ruleToDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *item)
	}

	return out, nil
}

func (r *scheduledRuleRepository) MarkRuleRan(ctx context.Context, ruleID uint, lastRunAt, nextRunAt time.Time) error {
	if err := r.db.Conn(ctx).Model(&models.ScheduledNotificationRule{}).
		Where("id = ?", ruleID).
		Updates(map[string]any{
			"last_run_at": lastRunAt,
			"next_run_at": nextRunAt,
		}).Error; err != nil {
		r.logger.Error().Err(err).Uint("rule_id", ruleID).Msg("Error marking scheduled rule as ran")
		return err
	}
	return nil
}

func ruleToModel(rule *entities.ScheduledRule) (*models.ScheduledNotificationRule, error) {
	params, err := json.Marshal(rule.SegmentParams)
	if err != nil {
		return nil, err
	}

	return &models.ScheduledNotificationRule{
		BusinessID:         rule.BusinessID,
		IntegrationID:      rule.IntegrationID,
		NotificationTypeID: rule.NotificationTypeID,
		WhatsappTemplateID: rule.WhatsappTemplateID,
		Name:               rule.Name,
		Description:        rule.Description,
		SegmentType:        rule.SegmentType,
		SegmentParams:      datatypes.JSON(params),
		Timezone:           rule.Timezone,
		FrequencyMinutes:   rule.FrequencyMinutes,
		SendWindowStart:    rule.SendWindowStart,
		SendWindowEnd:      rule.SendWindowEnd,
		CooldownDays:       rule.CooldownDays,
		DailySendCap:       rule.DailySendCap,
		BatchSizeCap:       rule.BatchSizeCap,
		RequiresOptIn:      rule.RequiresOptIn,
		Enabled:            rule.Enabled,
		LastRunAt:          rule.LastRunAt,
		NextRunAt:          rule.NextRunAt,
		CreatedByID:        rule.CreatedByID,
	}, nil
}

func ruleToDomain(model *models.ScheduledNotificationRule) (*entities.ScheduledRule, error) {
	rule := &entities.ScheduledRule{
		ID:                 model.ID,
		BusinessID:         model.BusinessID,
		IntegrationID:      model.IntegrationID,
		NotificationTypeID: model.NotificationTypeID,
		WhatsappTemplateID: model.WhatsappTemplateID,
		Name:               model.Name,
		Description:        model.Description,
		SegmentType:        model.SegmentType,
		Timezone:           model.Timezone,
		FrequencyMinutes:   model.FrequencyMinutes,
		SendWindowStart:    model.SendWindowStart,
		SendWindowEnd:      model.SendWindowEnd,
		CooldownDays:       model.CooldownDays,
		DailySendCap:       model.DailySendCap,
		BatchSizeCap:       model.BatchSizeCap,
		RequiresOptIn:      model.RequiresOptIn,
		Enabled:            model.Enabled,
		LastRunAt:          model.LastRunAt,
		NextRunAt:          model.NextRunAt,
		CreatedByID:        model.CreatedByID,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}

	if len(model.SegmentParams) > 0 {
		if err := json.Unmarshal(model.SegmentParams, &rule.SegmentParams); err != nil {
			return nil, err
		}
	}

	return rule, nil
}
