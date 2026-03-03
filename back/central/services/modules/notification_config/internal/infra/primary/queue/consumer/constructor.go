package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// DeliveryResultConsumer consume resultados de entrega de notificaciones
type DeliveryResultConsumer struct {
	rabbitMQ        rabbitmq.IQueue
	deliveryLogRepo ports.IDeliveryLogRepository
	logger          log.ILogger
}

// New crea un nuevo consumer de resultados de entrega
func New(rabbitMQ rabbitmq.IQueue, deliveryLogRepo ports.IDeliveryLogRepository, logger log.ILogger) *DeliveryResultConsumer {
	return &DeliveryResultConsumer{
		rabbitMQ:        rabbitMQ,
		deliveryLogRepo: deliveryLogRepo,
		logger:          logger,
	}
}
