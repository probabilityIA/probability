package softpymes

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/client"
	softpymescore "github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New crea una nueva instancia del módulo Softpymes
func New(
	config env.IConfig,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *softpymescore.SoftpymesCore {
	logger = logger.WithModule("softpymes")

	// 1. Cliente HTTP de Softpymes
	// Sin URL quemada: la URL siempre viene de integration_types.base_url / base_url_test.
	// El consumer valida que la URL no esté vacía antes de llamar al cliente.
	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Softpymes HTTP client initialized (URL dinámica desde integration_types)")

	// 2. Response Publisher (RabbitMQ)
	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Softpymes response publisher initialized")

	// 3. Invoice Request Consumer (escucha "invoicing.softpymes.requests")
	if rabbit != nil {
		invoiceRequestConsumer := consumer.NewInvoiceRequestConsumer(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ Softpymes invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting Softpymes invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ Softpymes invoice request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de facturación (Softpymes) deshabilitado")
	}

	// 4. Use Case
	useCase := app.New(httpClient, logger)
	logger.Info(context.Background()).Msg("✅ Softpymes use case initialized")

	logger.Info(context.Background()).Msg("✅ Softpymes bundle initialized")

	return softpymescore.New(useCase)
}
