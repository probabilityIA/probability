package consumerauthotp

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IConsumerAuthOTP interface {
	Start(ctx context.Context) error
}

type consumerAuthOTP struct {
	queue            rabbitmq.IQueue
	wa               ports.IWhatsApp
	credentialsCache ports.ICredentialsCache
	log              log.ILogger
}

func New(
	queue rabbitmq.IQueue,
	wa ports.IWhatsApp,
	credentialsCache ports.ICredentialsCache,
	logger log.ILogger,
) IConsumerAuthOTP {
	return &consumerAuthOTP{
		queue:            queue,
		wa:               wa,
		credentialsCache: credentialsCache,
		log:              logger,
	}
}
