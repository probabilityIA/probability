package siigo

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client"
	siigocore "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/queue"
	siigorepo "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	router *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	coreIntegration core.IIntegrationService,
) *siigocore.SiigoCore {
	logger = logger.WithModule("siigo")

	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ Siigo HTTP client initialized")

	responsePublisher := queue.New(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ Siigo response publisher initialized")

	var productRepo ports.IProductReadRepository
	if database != nil {
		productRepo = siigorepo.NewProductRepository(database, logger)
	}

	useCase := app.New(httpClient, coreIntegration, rabbit, productRepo, logger)
	logger.Info(context.Background()).Msg("✅ Siigo use case initialized")

	if router != nil && database != nil {
		webhookLogRepo := siigorepo.New(database, logger)
		webhookHandler := handlers.New(webhookLogRepo, coreIntegration, rabbit, logger)
		webhookHandler.RegisterRoutes(router)

		productHandler := handlers.NewProductHandler(useCase, logger)
		productHandler.RegisterRoutes(router)
		logger.Info(context.Background()).Msg("✅ Siigo webhook receiver and product sync registered")
	}

	if rabbit != nil {
		invoiceRequestConsumer := consumer.New(
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
			Msg("❌ RabbitMQ no disponible, consumer de facturacion (Siigo) deshabilitado")
	}

	logger.Info(context.Background()).Msg("✅ Siigo bundle initialized")

	return siigocore.New(useCase)
}
