package consumeralert

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IConsumerAlert define la interfaz del consumer de alertas de monitoreo
type IConsumerAlert interface {
	Start(ctx context.Context) error
}

// consumerAlert contiene las dependencias del consumer
type consumerAlert struct {
	queue            rabbitmq.IQueue
	wa               ports.IWhatsApp
	credentialsCache ports.ICredentialsCache
	log              log.ILogger
	phones           []string
}

// New crea una nueva instancia del consumer de alertas de monitoreo
func New(
	queue rabbitmq.IQueue,
	wa ports.IWhatsApp,
	credentialsCache ports.ICredentialsCache,
	logger log.ILogger,
	config env.IConfig,
) IConsumerAlert {
	raw := ""
	if config != nil {
		raw = config.Get(envAlertPhones)
	}
	phones := resolveAlertPhones(raw)

	logger.Info().
		Int("destinatarios", len(phones)).
		Strs("phones", phones).
		Msg("[AlertConsumer] Destinatarios de alertas configurados")

	return &consumerAlert{
		queue:            queue,
		wa:               wa,
		credentialsCache: credentialsCache,
		log:              logger,
		phones:           phones,
	}
}
