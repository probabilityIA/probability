package bold

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/app"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/primary/consumer"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/secondary/client"
	boldqueue "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	router *gin.RouterGroup,
	coreSvc core.IIntegrationCore,
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
) {
	logger = logger.WithModule("bold")
	ctx := context.Background()

	boldClient := client.New(logger)
	integrationRepo := repository.New(coreSvc, logger)
	responsePublisher := boldqueue.New(rabbit, logger)
	webhookPublisher := boldqueue.NewWebhookPublisher(rabbit, logger)

	useCase := app.New(boldClient, integrationRepo, responsePublisher, logger)
	webhookUseCase := app.NewWebhookUseCase(integrationRepo, webhookPublisher, logger)

	webhookHandlers := handlers.NewWebhookHandlers(webhookUseCase, logger)
	webhookHandlers.RegisterRoutes(router)

	if rabbit != nil {
		boldConsumer := consumer.New(rabbit, useCase, logger)
		go func() {
			if err := boldConsumer.Start(ctx); err != nil {
				logger.Error(ctx).Err(err).Msg("Bold consumer failed")
			}
		}()
		logger.Info(ctx).Msg("Bold consumer started")
	} else {
		logger.Warn(ctx).Msg("RabbitMQ not available - Bold consumer disabled")
	}
}
