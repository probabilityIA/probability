package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm/clause"
)

type campaignSendRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *campaignSendRepository) BulkCreateSends(ctx context.Context, sends []entities.CampaignSend) (int, error) {
	if len(sends) == 0 {
		return 0, nil
	}

	rows := make([]models.WhatsappCampaignSend, 0, len(sends))
	for _, send := range sends {
		rows = append(rows, models.WhatsappCampaignSend{
			CampaignID: send.CampaignID,
			ClientID:   send.ClientID,
			BusinessID: send.BusinessID,
			Phone:      send.Phone,
			Status:     models.CampaignSendStatusPending,
		})
	}

	result := r.db.Conn(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "campaign_id"}, {Name: "client_id"}},
			DoNothing: true,
		}).
		CreateInBatches(&rows, 200)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Error bulk creating campaign sends")
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (r *campaignSendRepository) ListPendingSends(ctx context.Context, campaignID uint, limit int) ([]entities.CampaignSend, error) {
	var rows []models.WhatsappCampaignSend

	if err := r.db.Conn(ctx).
		Where("campaign_id = ? AND status = ?", campaignID, models.CampaignSendStatusPending).
		Order("id ASC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		r.logger.Error().Err(err).Uint("campaign_id", campaignID).Msg("Error listing pending campaign sends")
		return nil, err
	}

	return sendsToDomain(rows), nil
}

type campaignSendRow struct {
	models.WhatsappCampaignSend
	ClientName string
}

func (r *campaignSendRepository) ListSends(ctx context.Context, campaignID uint, status string, page, pageSize int) ([]entities.CampaignSend, int64, error) {
	var total int64

	countQuery := r.db.Conn(ctx).Model(&models.WhatsappCampaignSend{}).Where("campaign_id = ?", campaignID)
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error counting campaign sends")
		return nil, 0, err
	}

	query := r.db.Conn(ctx).Model(&models.WhatsappCampaignSend{}).
		Select("whatsapp_campaign_sends.*, COALESCE(client.name, '') AS client_name").
		Joins("LEFT JOIN client ON client.id = whatsapp_campaign_sends.client_id").
		Where("whatsapp_campaign_sends.campaign_id = ?", campaignID)
	if status != "" {
		query = query.Where("whatsapp_campaign_sends.status = ?", status)
	}

	var rows []campaignSendRow
	if err := query.Order("whatsapp_campaign_sends.id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error listing campaign sends")
		return nil, 0, err
	}

	plain := make([]models.WhatsappCampaignSend, 0, len(rows))
	for i := range rows {
		plain = append(plain, rows[i].WhatsappCampaignSend)
	}

	out := sendsToDomain(plain)
	for i := range out {
		out[i].ClientName = rows[i].ClientName
	}

	return out, total, nil
}

func (r *campaignSendRepository) MarkSendQueued(ctx context.Context, sendID uint, queuedAt time.Time) error {
	if err := r.db.Conn(ctx).Model(&models.WhatsappCampaignSend{}).
		Where("id = ?", sendID).
		Updates(map[string]any{
			"status":    models.CampaignSendStatusQueued,
			"queued_at": queuedAt,
		}).Error; err != nil {
		r.logger.Error().Err(err).Uint("send_id", sendID).Msg("Error marking campaign send as queued")
		return err
	}
	return nil
}

func (r *campaignSendRepository) MarkSendResult(ctx context.Context, sendID uint, status, messageID, errorMessage string) error {
	updates := map[string]any{
		"status":        status,
		"message_id":    messageID,
		"error_message": errorMessage,
	}

	now := time.Now()
	switch status {
	case models.CampaignSendStatusSent:
		updates["sent_at"] = now
	case models.CampaignSendStatusDelivered:
		updates["delivered_at"] = now
	case models.CampaignSendStatusRead:
		updates["read_at"] = now
	case models.CampaignSendStatusReplied:
		updates["replied_at"] = now
	}

	if err := r.db.Conn(ctx).Model(&models.WhatsappCampaignSend{}).
		Where("id = ?", sendID).
		Updates(updates).Error; err != nil {
		r.logger.Error().Err(err).Uint("send_id", sendID).Msg("Error marking campaign send result")
		return err
	}
	return nil
}

func (r *campaignSendRepository) CountSentSince(ctx context.Context, campaignID uint, since time.Time) (int64, error) {
	var count int64

	if err := r.db.Conn(ctx).Model(&models.WhatsappCampaignSend{}).
		Where("campaign_id = ?", campaignID).
		Where("status IN ?", []string{
			models.CampaignSendStatusQueued,
			models.CampaignSendStatusSent,
			models.CampaignSendStatusDelivered,
			models.CampaignSendStatusRead,
			models.CampaignSendStatusReplied,
		}).
		Where("queued_at >= ?", since).
		Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Uint("campaign_id", campaignID).Msg("Error counting campaign sends since")
		return 0, err
	}

	return count, nil
}

func sendsToDomain(rows []models.WhatsappCampaignSend) []entities.CampaignSend {
	out := make([]entities.CampaignSend, 0, len(rows))
	for i := range rows {
		row := rows[i]

		var conversationID *string
		if row.ConversationID != nil {
			value := row.ConversationID.String()
			conversationID = &value
		}

		out = append(out, entities.CampaignSend{
			ID:             row.ID,
			CampaignID:     row.CampaignID,
			ClientID:       row.ClientID,
			BusinessID:     row.BusinessID,
			Phone:          row.Phone,
			ConversationID: conversationID,
			Status:         row.Status,
			MessageID:      row.MessageID,
			ErrorMessage:   row.ErrorMessage,
			QueuedAt:       row.QueuedAt,
			SentAt:         row.SentAt,
			DeliveredAt:    row.DeliveredAt,
			ReadAt:         row.ReadAt,
			RepliedAt:      row.RepliedAt,
		})
	}
	return out
}
