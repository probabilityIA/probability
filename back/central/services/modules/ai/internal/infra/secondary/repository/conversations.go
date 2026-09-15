package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) SaveMessage(ctx context.Context, record entities.MessageRecord) error {
	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		conversationID, err := ensureConversation(tx, record)
		if err != nil {
			return err
		}
		record.ConversationID = conversationID

		message := toMessageModel(record)
		created := tx.Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).Create(&message)
		if created.Error != nil {
			return created.Error
		}
		if created.RowsAffected == 0 {
			return nil
		}

		return tx.Model(&models.AssistantConversation{}).
			Where("id = ?", conversationID).
			Updates(map[string]any{
				"message_count":   gorm.Expr("message_count + 1"),
				"last_message_at": record.CreatedAt,
				"updated_at":      time.Now(),
			}).Error
	})
}

func ensureConversation(tx *gorm.DB, record entities.MessageRecord) (string, error) {
	var existing models.AssistantConversation
	err := tx.Select("id", "user_id").Where("id = ?", record.ConversationID).Take(&existing).Error
	switch {
	case err == nil && existing.UserID == record.UserID:
		return record.ConversationID, nil
	case err == nil:
		record.ConversationID = uuid.NewString()
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return "", err
	}

	conversation := models.AssistantConversation{
		ID:            record.ConversationID,
		BusinessID:    record.BusinessID,
		UserID:        record.UserID,
		StartedAt:     record.CreatedAt,
		LastMessageAt: record.CreatedAt,
	}
	if err := tx.Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).Create(&conversation).Error; err != nil {
		return "", err
	}
	return record.ConversationID, nil
}

func (r *Repository) ApplyFeedback(ctx context.Context, userID uint, messageID string, value int, at time.Time) (bool, error) {
	result := r.db.Conn(ctx).Model(&models.AssistantMessage{}).
		Where("id = ? AND user_id = ?", messageID, userID).
		Updates(map[string]any{"feedback": value, "feedback_at": at})
	return result.RowsAffected > 0, result.Error
}

func (r *Repository) ApplyClick(ctx context.Context, userID uint, messageID string, at time.Time) (bool, error) {
	result := r.db.Conn(ctx).Model(&models.AssistantMessage{}).
		Where("id = ? AND user_id = ?", messageID, userID).
		Update("clicked_at", gorm.Expr("COALESCE(clicked_at, ?)", at))
	return result.RowsAffected > 0, result.Error
}

func (r *Repository) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	conn := r.db.Conn(ctx)
	messages := conn.Where("created_at < ?", cutoff).Delete(&models.AssistantMessage{})
	if messages.Error != nil {
		return 0, messages.Error
	}
	if err := conn.Where("last_message_at < ?", cutoff).Delete(&models.AssistantConversation{}).Error; err != nil {
		return messages.RowsAffected, err
	}
	return messages.RowsAffected, nil
}

func toMessageModel(record entities.MessageRecord) models.AssistantMessage {
	return models.AssistantMessage{
		ID:               record.ID,
		ConversationID:   record.ConversationID,
		BusinessID:       record.BusinessID,
		UserID:           record.UserID,
		Pathname:         record.Pathname,
		Question:         record.Question,
		Answer:           record.Answer,
		DestinationKey:   record.DestinationKey,
		DestinationRoute: record.DestinationRoute,
		ErrorCode:        record.ErrorCode,
		Model:            record.Model,
		InputTokens:      record.InputTokens,
		OutputTokens:     record.OutputTokens,
		LatencyMs:        record.LatencyMs,
		CreatedAt:        record.CreatedAt,
	}
}
