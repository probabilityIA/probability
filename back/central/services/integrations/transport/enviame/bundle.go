package enviame

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New creates and initializes the Enviame transport provider.
// credentialResolver is used by the consumer to decrypt per-business API keys.
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	credentialResolver consumer.ICredentialResolver,
) {
	logger = logger.WithModule("transport.enviame")

	// 1. HTTP Client
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Enviame HTTP client initialized")

	// 2. Response Publisher
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Enviame response publisher initialized")

	// 3. Use Case
	useCase := app.New(httpClient, logger)
	logger.Info(context.Background()).Msg("✅ Enviame use case initialized")

	// 4. Request Consumer
	if rabbit != nil {
		requestConsumer := consumer.NewTransportRequestConsumer(
			rabbit,
			useCase,
			responsePublisher,
			credentialResolver,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Enviame transport request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Enviame transport request consumer in background...")
			if err := requestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Enviame transport request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de transporte (Enviame) deshabilitado")
	}

	logger.Info(context.Background()).Msg("✅ Enviame bundle initialized")
}
