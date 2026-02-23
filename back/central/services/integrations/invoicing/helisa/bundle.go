package helisa

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/infra/secondary/client"
	helisacore "github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea una nueva instancia del módulo Helisa
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *helisacore.HelisaCore {
	logger = logger.WithModule("helisa")

	// 1. Cliente HTTP de Helisa
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Helisa HTTP client initialized")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Helisa response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.helisa.requests")
	if rabbit != nil {
		invoiceRequestConsumer := consumer.NewInvoiceRequestConsumer(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Helisa invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Helisa invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Helisa invoice request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Helisa) deshabilitado")
	}

	// 4. Use Case — contiene toda la lógica de negocio
	useCase := app.New(httpClient, coreIntegration, logger)
	logger.Info(context.Background()).Msg("✅ Helisa use case initialized")

	logger.Info(context.Background()).Msg("✅ Helisa bundle initialized")

	return helisacore.New(useCase)
}
