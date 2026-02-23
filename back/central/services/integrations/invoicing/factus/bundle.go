package factus

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client"
	factuscore "github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea una nueva instancia del módulo Factus
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *factuscore.FactusCore {
	logger = logger.WithModule("factus")

	// 1. Cliente HTTP de Factus (adapter secundario)
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Factus HTTP client initialized")

	// 2. Use Case — contiene toda la lógica de negocio
	useCase := app.New(httpClient, coreIntegration, logger)
	logger.Info(context.Background()).Msg("✅ Factus use case initialized")

	// 3. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Factus response publisher initialized")

	// 4. Invoice Request Consumer (escucha "invoicing.factus.requests")
	if rabbit != nil {
		invoiceRequestConsumer := consumer.NewInvoiceRequestConsumer(
			rabbit,
			useCase,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Factus invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Factus invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Factus invoice request consumer failed to start or stopped with error")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Factus) deshabilitado")
	}

	logger.Info(context.Background()).Msg("✅ Factus bundle initialized (HTTP client + RabbitMQ async consumer)")

	return factuscore.New(useCase)
}
