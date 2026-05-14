package whatsapp_persistence_consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type PersistenceConsumer struct {
	rabbitMQ  rabbitmq.IQueue
	persister ports.IWhatsAppPersister
	logger    log.ILogger
}

func New(rabbitMQ rabbitmq.IQueue, persister ports.IWhatsAppPersister, logger log.ILogger) *PersistenceConsumer {
	return &PersistenceConsumer{
		rabbitMQ:  rabbitMQ,
		persister: persister,
		logger:    logger,
	}
}
