package siigo

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client"
	siigocore "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea una nueva instancia del módulo Siigo
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *siigocore.SiigoCore {
	logger = logger.WithModule("siigo")

	// 1. Cliente HTTP de Siigo
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Siigo HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Siigo response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.siigo.requests")
	if rabbit != nil {
		invoiceRequestConsumer := consumer.NewInvoiceRequestConsumer(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Siigo invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Siigo invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Siigo invoice request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Siigo) deshabilitado")
	}

	// 4. Use Case — contiene toda la lógica de negocio
	useCase := app.New(httpClient, coreIntegration, logger)
	logger.Info(context.Background()).Msg("✅ Siigo use case initialized")

	logger.Info(context.Background()).Msg("✅ Siigo bundle initialized")

	return siigocore.New(useCase)
}
