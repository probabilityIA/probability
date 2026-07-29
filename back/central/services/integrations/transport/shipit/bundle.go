package shipit

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
	credentialResolver consumer.ICredentialResolver,
) {
	logger = logger.WithModule("transport.shipit")

	httpClient := client.New(logger)
	logger.Info(context.Background()).Msg("Shipit HTTP client initialized")

	responsePublisher := queue.NewResponsePublisher(rabbit, logger)
	webhookPublisher := queue.NewWebhookResponsePublisher(responsePublisher)

	useCase := app.New(httpClient, logger)
	logger.Info(context.Background()).Msg("Shipit use case initialized")

	if rabbit != nil {
		requestConsumer := consumer.NewTransportRequestConsumer(
			rabbit,
			useCase,
			responsePublisher,
			webhookPublisher,
			credentialResolver,
			logger,
		)

		go func() {
			ctx := context.Background()
			logger.Info(ctx).Msg("Starting Shipit transport request consumer in background...")
			if err := requestConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("Shipit transport request consumer failed")
			}
		}()
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ no disponible, consumer de transporte (Shipit) deshabilitado")
	}

	webhookUC := app.NewWebhookUseCase(webhookPublisher, logger)
	webhookHandlers := handlers.New(webhookUC, logger, rabbit)
	webhookHandlers.RegisterRoutes(router)

	if rabbit != nil {
		webhookQueueConsumer := consumer.NewWebhookQueueConsumer(rabbit, webhookUC, logger)
		webhookQueueConsumer.Start(context.Background())
	}

	logger.Info(context.Background()).Msg("Shipit bundle initialized")
}
