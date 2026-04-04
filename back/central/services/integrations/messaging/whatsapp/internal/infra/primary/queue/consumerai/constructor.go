package consumerai

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IConsumer define la interfaz del consumer de respuestas AI
type IConsumer interface {
	Start(ctx context.Context) error
}

// consumer contiene las dependencias del consumer
type consumer struct {
	queue            rabbitmq.IQueue
	wa               ports.IWhatsApp
	credentialsCache ports.ICredentialsCache
	log              log.ILogger
}

// New crea una nueva instancia del consumer de respuestas AI
func New(
	queue rabbitmq.IQueue,
	wa ports.IWhatsApp,
	credentialsCache ports.ICredentialsCache,
	logger log.ILogger,
) IConsumer {
	return &consumer{
		queue:            queue,
		wa:               wa,
		credentialsCache: credentialsCache,
		log:              logger.WithModule("whatsapp-ai-response-consumer"),
	}
}
