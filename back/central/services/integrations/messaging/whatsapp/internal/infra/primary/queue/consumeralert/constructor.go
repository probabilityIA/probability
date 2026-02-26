package consumeralert

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// IConsumerAlert define la interfaz del consumer de alertas de monitoreo
type IConsumerAlert interface {
	Start(ctx context.Context) error
}

// consumerAlert contiene las dependencias del consumer
type consumerAlert struct {
	queue           rabbitmq.IQueue
	wa              ports.IWhatsApp
	integrationRepo ports.IIntegrationRepository
	log             log.ILogger
}

// New crea una nueva instancia del consumer de alertas de monitoreo
func New(
	queue rabbitmq.IQueue,
	wa ports.IWhatsApp,
	integrationRepo ports.IIntegrationRepository,
	logger log.ILogger,
) IConsumerAlert {
	return &consumerAlert{
		queue:           queue,
		wa:              wa,
		integrationRepo: integrationRepo,
		log:             logger,
	}
}
