package worldoffice

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/infra/secondary/client"
	worldofficecore "github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea una nueva instancia del módulo World Office
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *worldofficecore.WorldOfficeCore {
	logger = logger.WithModule("world_office")

	// 1. Cliente HTTP de World Office
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ World Office HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ World Office response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.world_office.requests")
	if rabbit != nil {
		invoiceRequestConsumer := consumer.NewInvoiceRequestConsumer(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ World Office invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting World Office invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ World Office invoice request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (World Office) deshabilitado")
	}

	// 4. Use Case — contiene toda la lógica de negocio
	useCase := app.New(httpClient, coreIntegration, logger)
	logger.Info(context.Background()).Msg("✅ World Office use case initialized")

	logger.Info(context.Background()).Msg("✅ World Office bundle initialized")

	return worldofficecore.New(useCase)
}
