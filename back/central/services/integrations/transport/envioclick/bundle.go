package envioclick

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	router *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	credentialResolver consumer.ICredentialResolver,
) {
	logger = logger.WithModule("transport.envioclick")

	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("✅ EnvioClick HTTP client initialized")

	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	logger.Info(context.Background()).Msg("✅ EnvioClick response publisher initialized")

	useCase := app.New(httpClient, logger)
	logger.Info(context.Background()).Msg("✅ EnvioClick use case initialized")

	webhookRepo := repository.New(database)
	webhookPublisher := queue.NewWebhookResponsePublisher(responsePublisher)
	syncUseCase := app.NewSyncUseCase(httpClient, webhookRepo, webhookPublisher, logger)

	if rabbit != nil {
		requestConsumer := consumer.NewTransportRequestConsumer(
			rabbit,
			useCase,
			syncUseCase,
			responsePublisher,
			credentialResolver,
			logger,
		)
		logger.Info(context.Background()).Msg("✅ EnvioClick transport request consumer initialized")

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("🚀 Starting EnvioClick transport request consumer in background...")
			if err := requestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("❌ EnvioClick transport request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("❌ RabbitMQ no disponible, consumer de transporte (EnvioClick) deshabilitado")
	}

	webhookUC := app.NewWebhookUseCase(webhookRepo, webhookPublisher, logger)
	webhookHandlers := handlers.New(webhookUC, logger)
	webhookHandlers.RegisterRoutes(router)
	logger.Info(context.Background()).Msg("✅ EnvioClick webhook handler registered")

	logger.Info(context.Background()).Msg("✅ EnvioClick bundle initialized")
}
