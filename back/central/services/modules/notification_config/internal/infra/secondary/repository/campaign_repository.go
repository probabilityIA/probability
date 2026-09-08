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

type campaignRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *campaignRepository) CreateCampaign(ctx context.Context, campaign *entities.Campaign) error {
	model, err := campaignToModel(campaign)
	if err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.logger.Error().Err(err).Str("name", campaign.Name).Msg("Error creating campaign")
		return err
	}

	campaign.ID = model.ID
	campaign.CreatedAt = model.CreatedAt
	campaign.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *campaignRepository) UpdateCampaign(ctx context.Context, campaign *entities.Campaign) error {
	model, err := campaignToModel(campaign)
	if err != nil {
		return err
	}

	if err := r.db.Conn(ctx).Model(&models.WhatsappCampaign{}).
		Where("id = ? AND business_id = ?", campaign.ID, campaign.BusinessID).
		Updates(map[string]any{
			"integration_id":       model.IntegrationID,
			"whatsapp_template_id": model.WhatsappTemplateID,
			"name":                 model.Name,
			"description":          model.Description,
			"sender_name":          model.SenderName,
			"audience_type":        model.AudienceType,
			"audience_params":      model.AudienceParams,
			"variable_values":      model.VariableValues,
			"timezone":             model.Timezone,
			"send_window_start":    model.SendWindowStart,
			"send_window_end":      model.SendWindowEnd,
			"scheduled_at":         model.ScheduledAt,
			"daily_send_cap":       model.DailySendCap,
			"batch_size":           model.BatchSize,
			"status":               model.Status,
			"audience_count":       model.AudienceCount,
			"started_at":           model.StartedAt,
			"finished_at":          model.FinishedAt,
			"error_message":        model.ErrorMessage,
		}).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", campaign.ID).Msg("Error updating campaign")
		return err
	}

	return nil
}

func (r *campaignRepository) GetCampaignByID(ctx context.Context, id uint) (*entities.Campaign, error) {
	var model models.WhatsappCampaign

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("id", id).Msg("Error getting campaign")
		return nil, err
	}

	return campaignToDomain(&model)
}

func (r *campaignRepository) ListCampaigns(ctx context.Context, businessID uint, status string, page, pageSize int) ([]entities.Campaign, int64, error) {
	var rows []models.WhatsappCampaign
	var total int64

	query := r.db.Conn(ctx).Model(&models.WhatsappCampaign{}).Where("business_id = ?", businessID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error counting campaigns")
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing campaigns")
		return nil, 0, err
	}

	out := make([]entities.Campaign, 0, len(rows))
	for i := range rows {
		item, err := campaignToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *item)
	}

	return out, total, nil
}

func (r *campaignRepository) DeleteCampaign(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.WhatsappCampaign{}, id).Error; err != nil {
		r.logger.Error().Err(err).Uint("id", id).Msg("Error deleting campaign")
		return err
	}
	return nil
}

func (r *campaignRepository) ListDueCampaigns(ctx context.Context, now time.Time, limit int) ([]entities.Campaign, error) {
	var rows []models.WhatsappCampaign

	if err := r.db.Conn(ctx).
		Where("status IN ?", []string{models.CampaignStatusScheduled, models.CampaignStatusRunning}).
		Where("scheduled_at IS NULL OR scheduled_at <= ?", now).
		Order("scheduled_at ASC NULLS FIRST").
		Limit(limit).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing due campaigns")
		return nil, err
	}

	out := make([]entities.Campaign, 0, len(rows))
	for i := range rows {
		item, err := campaignToDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *item)
	}

	return out, nil
}

func (r *campaignRepository) UpdateCampaignCounters(ctx context.Context, campaignID uint) error {
	if err := r.db.Conn(ctx).Exec(`
		UPDATE whatsapp_campaigns c SET
			queued_count = s.queued,
			sent_count = s.sent,
			failed_count = s.failed,
			skipped_count = s.skipped,
			replied_count = s.replied,
			updated_at = NOW()
		FROM (
			SELECT
				COUNT(*) FILTER (WHERE status = 'queued') AS queued,
				COUNT(*) FILTER (WHERE status IN ('sent','delivered','read','replied')) AS sent,
				COUNT(*) FILTER (WHERE status = 'failed') AS failed,
				COUNT(*) FILTER (WHERE status = 'skipped') AS skipped,
				COUNT(*) FILTER (WHERE status = 'replied') AS replied
			FROM whatsapp_campaign_sends
			WHERE campaign_id = ? AND deleted_at IS NULL
		) s
		WHERE c.id = ?`, campaignID, campaignID).Error; err != nil {
		r.logger.Error().Err(err).Uint("campaign_id", campaignID).Msg("Error updating campaign counters")
		return err
	}
	return nil
}

func (r *campaignRepository) MarkCampaignBatch(ctx context.Context, campaignID uint, lastBatchAt time.Time) error {
	if err := r.db.Conn(ctx).Model(&models.WhatsappCampaign{}).
		Where("id = ?", campaignID).
		Update("last_batch_at", lastBatchAt).Error; err != nil {
		r.logger.Error().Err(err).Uint("campaign_id", campaignID).Msg("Error marking campaign batch")
		return err
	}
	return nil
}

func campaignToModel(campaign *entities.Campaign) (*models.WhatsappCampaign, error) {
	params, err := json.Marshal(campaign.AudienceParams)
	if err != nil {
		return nil, err
	}

	values, err := json.Marshal(campaign.VariableValues)
	if err != nil {
		return nil, err
	}

	return &models.WhatsappCampaign{
		BusinessID:         campaign.BusinessID,
		IntegrationID:      campaign.IntegrationID,
		WhatsappTemplateID: campaign.WhatsappTemplateID,
		Name:               campaign.Name,
		Description:        campaign.Description,
		SenderName:         campaign.SenderName,
		AudienceType:       campaign.AudienceType,
		AudienceParams:     datatypes.JSON(params),
		VariableValues:     datatypes.JSON(values),
		Timezone:           campaign.Timezone,
		SendWindowStart:    campaign.SendWindowStart,
		SendWindowEnd:      campaign.SendWindowEnd,
		ScheduledAt:        campaign.ScheduledAt,
		DailySendCap:       campaign.DailySendCap,
		BatchSize:          campaign.BatchSize,
		Status:             campaign.Status,
		AudienceCount:      campaign.AudienceCount,
		StartedAt:          campaign.StartedAt,
		FinishedAt:         campaign.FinishedAt,
		ErrorMessage:       campaign.ErrorMessage,
		CreatedByID:        campaign.CreatedByID,
	}, nil
}

func campaignToDomain(model *models.WhatsappCampaign) (*entities.Campaign, error) {
	campaign := &entities.Campaign{
		ID:                 model.ID,
		BusinessID:         model.BusinessID,
		IntegrationID:      model.IntegrationID,
		WhatsappTemplateID: model.WhatsappTemplateID,
		Name:               model.Name,
		Description:        model.Description,
		SenderName:         model.SenderName,
		AudienceType:       model.AudienceType,
		Timezone:           model.Timezone,
		SendWindowStart:    model.SendWindowStart,
		SendWindowEnd:      model.SendWindowEnd,
		ScheduledAt:        model.ScheduledAt,
		DailySendCap:       model.DailySendCap,
		BatchSize:          model.BatchSize,
		Status:             model.Status,
		AudienceCount:      model.AudienceCount,
		QueuedCount:        model.QueuedCount,
		SentCount:          model.SentCount,
		FailedCount:        model.FailedCount,
		SkippedCount:       model.SkippedCount,
		RepliedCount:       model.RepliedCount,
		StartedAt:          model.StartedAt,
		FinishedAt:         model.FinishedAt,
		LastBatchAt:        model.LastBatchAt,
		ErrorMessage:       model.ErrorMessage,
		CreatedByID:        model.CreatedByID,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}

	if len(model.AudienceParams) > 0 {
		if err := json.Unmarshal(model.AudienceParams, &campaign.AudienceParams); err != nil {
			return nil, err
		}
	}

	if len(model.VariableValues) > 0 {
		if err := json.Unmarshal(model.VariableValues, &campaign.VariableValues); err != nil {
			return nil, err
		}
	}

	return campaign, nil
}
