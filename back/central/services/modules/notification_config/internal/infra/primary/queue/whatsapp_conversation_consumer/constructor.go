package whatsapp_conversation_consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ConversationConsumer consume eventos de conversaciones WhatsApp y los persiste en DB
type ConversationConsumer struct {
	rabbitMQ  rabbitmq.IQueue
	persister ports.IWhatsAppPersister
	logger    log.ILogger
}

// New crea un nuevo consumer de conversaciones WhatsApp
func New(rabbitMQ rabbitmq.IQueue, persister ports.IWhatsAppPersister, logger log.ILogger) *ConversationConsumer {
	return &ConversationConsumer{
		rabbitMQ:  rabbitMQ,
		persister: persister,
		logger:    logger,
	}
}
