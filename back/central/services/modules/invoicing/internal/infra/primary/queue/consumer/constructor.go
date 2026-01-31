package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Consumers agrupa todos los consumers del módulo
type Consumers struct {
	Order *OrderConsumer
	Retry *RetryConsumer
}

// NewConsumers crea todos los consumers del módulo de facturación
func NewConsumers(
	queue rabbitmq.IQueue,
	useCase ports.IUseCase,
	syncLogRepo ports.IInvoiceSyncLogRepository,
	logger log.ILogger,
) *Consumers {
	return &Consumers{
		Order: NewOrderConsumer(queue, useCase, logger),
		Retry: NewRetryConsumer(syncLogRepo, useCase, logger),
	}
}
