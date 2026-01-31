package repository

import (
		"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New crea todas las instancias de repositorios en un solo constructor
func New(
	database db.IDatabase,
	logger log.ILogger,
	encryptionKey []byte,
) (
	ports.IConversationRepository,
	ports.IMessageLogRepository,
	ports.IIntegrationRepository,
) {
	// 1. ConversationRepository
	conversationRepo := &ConversationRepository{
		db:  database,
		log: logger.WithModule("whatsapp-conversation-repo"),
	}

	// 2. MessageLogRepository
	messageLogRepo := &MessageLogRepository{
		db:  database,
		log: logger.WithModule("whatsapp-message-log-repo"),
	}

	// 3. IntegrationRepository
	integrationRepo := &IntegrationRepository{
		db:            database,
		log:           logger.WithModule("whatsapp-integration-repo"),
		encryptionKey: encryptionKey,
	}

	return conversationRepo, messageLogRepo, integrationRepo
}
