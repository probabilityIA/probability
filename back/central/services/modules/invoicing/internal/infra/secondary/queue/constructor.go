package queue

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea un nuevo event publisher
func New(queue rabbitmq.IQueue, logger log.ILogger) ports.IEventPublisher {
	return NewEventPublisher(queue, logger)
}
