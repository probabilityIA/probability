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

type scheduledSendRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func (r *scheduledSendRepository) ReserveSend(ctx context.Context, send *entities.ScheduledSend) (bool, error) {
	model := &models.ScheduledNotificationSend{
		RuleID:     send.RuleID,
		ClientID:   send.ClientID,
		SendDate:   send.SendDate,
		RunID:      send.RunID,
		BusinessID: send.BusinessID,
		Phone:      send.Phone,
		Status:     entities.SendStatusQueued,
		QueuedAt:   time.Now(),
	}

	result := r.db.Conn(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rule_id"}, {Name: "client_id"}, {Name: "send_date"}},
		DoNothing: true,
	}).Create(model)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).
			Uint("rule_id", send.RuleID).
			Uint("client_id", send.ClientID).
			Msg("Error reserving scheduled send")
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	send.ID = model.ID
	send.QueuedAt = model.QueuedAt
	send.Status = model.Status
	return true, nil
}

func (r *scheduledSendRepository) MarkSendResult(ctx context.Context, sendID uint, status, messageID, errorMessage string) error {
	updates := map[string]any{
		"status":        status,
		"message_id":    messageID,
		"error_message": errorMessage,
	}
	if status == entities.SendStatusSent {
		now := time.Now()
		updates["sent_at"] = &now
	}

	if err := r.db.Conn(ctx).Model(&models.ScheduledNotificationSend{}).
		Where("id = ?", sendID).
		Updates(updates).Error; err != nil {
		r.logger.Error().Err(err).Uint("send_id", sendID).Msg("Error marking scheduled send result")
		return err
	}

	return nil
}

func (r *scheduledSendRepository) CountSendsToday(ctx context.Context, ruleID uint, sendDate string) (int64, error) {
	var total int64

	if err := r.db.Conn(ctx).Model(&models.ScheduledNotificationSend{}).
		Where("rule_id = ? AND send_date = ? AND status <> ?", ruleID, sendDate, entities.SendStatusDiscarded).
		Count(&total).Error; err != nil {
		r.logger.Error().Err(err).Uint("rule_id", ruleID).Msg("Error counting scheduled sends for the day")
		return 0, err
	}

	return total, nil
}
