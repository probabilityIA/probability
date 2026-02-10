package consumer

import (
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Consumers agrupa todos los consumers del módulo
type Consumers struct {
<<<<<<< HEAD
	Order *OrderConsumer
	Retry *RetryConsumer
=======
	Order       *OrderConsumer
	Retry       *RetryConsumer
	BulkInvoice *BulkInvoiceConsumer
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}

// NewConsumers crea todos los consumers del módulo de facturación
func NewConsumers(
	queue rabbitmq.IQueue,
	useCase ports.IUseCase,
<<<<<<< HEAD
	syncLogRepo ports.IInvoiceSyncLogRepository,
	logger log.ILogger,
) *Consumers {
	return &Consumers{
		Order: NewOrderConsumer(queue, useCase, logger),
		Retry: NewRetryConsumer(syncLogRepo, useCase, logger),
=======
	repo ports.IRepository,
	ssePublisher ports.IInvoiceSSEPublisher,
	logger log.ILogger,
) *Consumers {
	return &Consumers{
		Order:       NewOrderConsumer(queue, useCase, logger),
		Retry:       NewRetryConsumer(repo, useCase, logger),
		BulkInvoice: NewBulkInvoiceConsumer(queue, useCase, repo, ssePublisher, logger),
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	}
}
