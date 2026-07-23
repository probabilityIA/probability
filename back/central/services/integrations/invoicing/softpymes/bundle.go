package softpymes

import (
	"context"
	"strconv"

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

const defaultInvoiceWorkers = 3

func New(
	config env.IConfig,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *softpymescore.SoftpymesCore {
	logger = logger.WithModule("softpymes")

	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("Softpymes HTTP client initialized (URL dinamica desde integration_types)")

	responsePublisher := queue.New(rabbit, logger)
	logger.Info(context.Background()).Msg("Softpymes response publisher initialized")

	if rabbit != nil {
		workers := defaultInvoiceWorkers
		if raw := config.Get("SOFTPYMES_INVOICE_WORKERS"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
				workers = parsed
			} else {
				logger.Warn(context.Background()).
					Str("value", raw).
					Int("default", defaultInvoiceWorkers).
					Msg("SOFTPYMES_INVOICE_WORKERS invalido, usando default")
			}
		}

		invoiceRequestConsumer := consumer.New(
			rabbit,
			coreIntegration,
			httpClient,
			responsePublisher,
			logger,
			workers,
		)
		logger.Info(context.Background()).
			Int("workers", workers).
			Msg("Softpymes invoice request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("Starting Softpymes invoice request consumer in background...")
			if err := invoiceRequestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("Softpymes invoice request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ no disponible, consumer de facturacion (Softpymes) deshabilitado")
	}

	useCase := app.New(httpClient, logger)
	logger.Info(context.Background()).Msg("Softpymes use case initialized")

	logger.Info(context.Background()).Msg("Softpymes bundle initialized")

	return softpymescore.New(useCase)
}
