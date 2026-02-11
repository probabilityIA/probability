package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Consumers agrupa todos los consumers del módulo
type Consumers struct {
	Order       *OrderConsumer
	Retry       *RetryConsumer
	BulkInvoice *BulkInvoiceConsumer
	Response    *ResponseConsumer // Consumer de responses de proveedores
}

// NewConsumers crea todos los consumers del módulo de facturación
func NewConsumers(
	queue rabbitmq.IQueue,
	useCase ports.IUseCase,
	repo ports.IRepository,
	ssePublisher ports.IInvoiceSSEPublisher,
	eventPublisher ports.IEventPublisher,
	logger log.ILogger,
) *Consumers {
	return &Consumers{
		Order:       NewOrderConsumer(queue, useCase, logger),
		Retry:       NewRetryConsumer(repo, useCase, logger),
		BulkInvoice: NewBulkInvoiceConsumer(queue, useCase, repo, ssePublisher, logger),
		Response:    NewResponseConsumer(queue, repo, ssePublisher, eventPublisher, logger),
	}
}
