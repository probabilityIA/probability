package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Consumers agrupa todos los consumers del módulo de pagos
type Consumers struct {
	Response *ResponseConsumer
	Retry    *RetryConsumer
}

// NewConsumers crea todos los consumers del módulo de pagos
func NewConsumers(
	queue rabbitmq.IQueue,
	useCase ports.IUseCase,
	repo ports.IRepository,
	ssePublisher ports.ISSEPublisher,
	logger log.ILogger,
) *Consumers {
	return &Consumers{
		Response: NewResponseConsumer(queue, useCase, logger),
		Retry:    NewRetryConsumer(repo, useCase, logger),
	}
}
