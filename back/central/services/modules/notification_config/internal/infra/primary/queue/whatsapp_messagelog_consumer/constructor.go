package whatsapp_messagelog_consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// MessageLogConsumer consume eventos de message logs WhatsApp y los persiste en DB
type MessageLogConsumer struct {
	rabbitMQ  rabbitmq.IQueue
	persister ports.IWhatsAppPersister
	logger    log.ILogger
}

// New crea un nuevo consumer de message logs WhatsApp
func New(rabbitMQ rabbitmq.IQueue, persister ports.IWhatsAppPersister, logger log.ILogger) *MessageLogConsumer {
	return &MessageLogConsumer{
		rabbitMQ:  rabbitMQ,
		persister: persister,
		logger:    logger,
	}
}
