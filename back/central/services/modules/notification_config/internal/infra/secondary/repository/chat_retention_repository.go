package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

const chatPurgeBatchSize = 5000

type chatRetentionRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

func NewChatRetentionRepository(database db.IDatabase, logger log.ILogger) ports.IChatRetentionRepository {
	return &chatRetentionRepository{db: database, logger: logger.WithModule("chat_retention_repository")}
}

func (r *chatRetentionRepository) PurgeChatHistory(ctx context.Context, cutoff time.Time) (entities.ChatPurgeResult, error) {
	var result entities.ChatPurgeResult

	for {
		deleted := r.db.Conn(ctx).Exec(`
			DELETE FROM whatsapp_message_logs
			WHERE id IN (
				SELECT id FROM whatsapp_message_logs
				WHERE created_at < ?
				LIMIT ?
			)`, cutoff, chatPurgeBatchSize)
		if deleted.Error != nil {
			return result, deleted.Error
		}
		result.Messages += deleted.RowsAffected
		if deleted.RowsAffected < chatPurgeBatchSize {
			break
		}
	}

	conversations := r.db.Conn(ctx).Exec(`
		DELETE FROM whatsapp_conversations c
		WHERE c.created_at < ?
		  AND NOT EXISTS (
			SELECT 1 FROM whatsapp_message_logs ml WHERE ml.conversation_id = c.id
		  )`, cutoff)
	if conversations.Error != nil {
		return result, conversations.Error
	}
	result.Conversations = conversations.RowsAffected

	reads := r.db.Conn(ctx).Exec(`
		DELETE FROM whatsapp_conversation_reads rd
		WHERE rd.updated_at < ?
		  AND NOT EXISTS (
			SELECT 1 FROM whatsapp_conversations c
			WHERE c.business_id = rd.business_id
			  AND regexp_replace(c.phone_number, '[^0-9]', '', 'g') = rd.phone_key
		  )`, cutoff)
	if reads.Error != nil {
		return result, reads.Error
	}
	result.Reads = reads.RowsAffected

	return result, nil
}
