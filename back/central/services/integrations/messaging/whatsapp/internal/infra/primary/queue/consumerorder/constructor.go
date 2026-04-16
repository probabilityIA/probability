package consumerorder

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IConsumer define la interfaz para el consumer de órdenes
type IConsumer interface {
	Start(ctx context.Context) error
}

// consumer contiene las dependencias del consumer
type consumer struct {
	queue   rabbitmq.IQueue
	useCase usecasemessaging.IUseCase
	log     log.ILogger
}

// New crea una nueva instancia del consumer de órdenes
func New(
	queue rabbitmq.IQueue,
	useCase usecasemessaging.IUseCase,
	logger log.ILogger,
) IConsumer {
	return &consumer{
		queue:   queue,
		useCase: useCase,
		log:     logger,
	}
}
