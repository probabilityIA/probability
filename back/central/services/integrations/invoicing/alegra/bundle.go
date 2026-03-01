package alegra

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/infra/secondary/client"
	alegracore "github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea una nueva instancia del módulo Alegra
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *alegracore.AlegraCore {
	logger = logger.WithModule("alegra")

	// 1. Cliente HTTP de Alegra
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Alegra HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.New(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Alegra response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.alegra.requests")
	if rabbit != nil {
		invoiceRequestConsumer := consumer.New(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Alegra invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Alegra invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Alegra invoice request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Alegra) deshabilitado")
	}

	// 4. Use Case — contiene toda la lógica de negocio
	useCase := app.New(httpClient, coreIntegration, logger)
	logger.Info(context.Background()).Msg("✅ Alegra use case initialized")

	logger.Info(context.Background()).Msg("✅ Alegra bundle initialized")

	return alegracore.New(useCase)
}
