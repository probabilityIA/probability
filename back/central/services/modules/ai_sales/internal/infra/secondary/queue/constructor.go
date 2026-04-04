package queue

import (
	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type responsePublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewResponsePublisher crea un publisher para enviar respuestas AI a WhatsApp
func NewResponsePublisher(rabbit rabbitmq.IQueue, logger log.ILogger) domain.IAIResponsePublisher {
	return &responsePublisher{
		rabbit: rabbit,
		log:    logger,
	}
}

type orderPublisher struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// NewOrderPublisher crea un publisher para enviar ordenes al queue canonical
func NewOrderPublisher(rabbit rabbitmq.IQueue, logger log.ILogger) domain.IAIOrderPublisher {
	return &orderPublisher{
		rabbit: rabbit,
		log:    logger,
	}
}
