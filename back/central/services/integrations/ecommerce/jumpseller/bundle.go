package jumpseller

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/primary/handlers"
	jumpsellerqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client"
	jumpsellercore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/queue"
	jumpsellerproductrepo "github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/repository"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
	database db.IDatabase,
	coreIntegration integrationcore.IIntegrationCore,
) integrationcore.IIntegrationContract {
	logger = logger.WithModule("jumpseller")

	httpClient := client.New()
	integrationService := jumpsellercore.NewIntegrationService(coreIntegration)
	productRepo := jumpsellerproductrepo.New(database, logger)

	var orderPublisher = queue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, Jumpseller orders will not be published to queue")
	}

	uc := usecases.New(httpClient, integrationService, orderPublisher, productRepo, rabbitMQ, logger)

	handler := handlers.New(uc, logger)
	handler.RegisterRoutes(router, logger)

	if rabbitMQ != nil {
		pushConsumer := jumpsellerqueue.NewInventoryPushConsumer(rabbitMQ, uc, logger)
		pushConsumer.Start(context.Background())
	}

	baseURL := config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = config.Get("URL_BASE_SWAGGER")
	}
	if baseURL != "" {
		coreIntegration.OnIntegrationCreated(integrationcore.IntegrationTypeJumpseller, func(obsCtx context.Context, integration *integrationcore.PublicIntegration) {
			go func() {
				bgCtx := context.Background()
				integrationID := fmt.Sprintf("%d", integration.ID)
				if err := uc.CreateWebhooks(bgCtx, integrationID, baseURL); err != nil {
					logger.Error(bgCtx).Err(err).Str("integration_id", integrationID).Msg("Error al crear webhooks automaticamente para Jumpseller")
				} else {
					logger.Info(bgCtx).Str("integration_id", integrationID).Msg("Webhooks creados automaticamente para Jumpseller")
				}
			}()
		})
	} else {
		logger.Warn(context.Background()).Msg("Ni WEBHOOK_BASE_URL ni URL_BASE_SWAGGER configuradas, no se crearan webhooks automaticamente para Jumpseller")
	}

	return jumpsellercore.New(uc)
}
