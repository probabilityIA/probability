package tiendanube

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/primary/handlers"
	tiendanubequeue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/client"
	tncore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/queue"
	tnproductrepo "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/infra/secondary/repository"
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
	logger = logger.WithModule("tiendanube")

	httpClient := client.New()
	integrationService := tncore.NewIntegrationService(coreIntegration)
	productRepo := tnproductrepo.New(database, logger)

	var orderPublisher = queue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, Tiendanube orders will not be published to queue")
	}

	uc := usecases.New(httpClient, integrationService, orderPublisher, productRepo, rabbitMQ, logger)

	handler := handlers.New(uc, coreIntegration, config, logger)
	handler.RegisterRoutes(router, logger)

	if rabbitMQ != nil {
		pushConsumer := tiendanubequeue.NewInventoryPushConsumer(rabbitMQ, uc, logger)
		pushConsumer.Start(context.Background())
	}

	baseURL := config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = config.Get("URL_BASE_SWAGGER")
	}
	if baseURL != "" {
		coreIntegration.OnIntegrationCreated(integrationcore.IntegrationTypeTiendanube, func(obsCtx context.Context, integration *integrationcore.PublicIntegration) {
			go func() {
				bgCtx := context.Background()
				integrationID := fmt.Sprintf("%d", integration.ID)
				if _, err := uc.CreateWebhooks(bgCtx, integrationID, baseURL); err != nil {
					logger.Error(bgCtx).Err(err).Str("integration_id", integrationID).Msg("Error al crear webhooks automaticamente para Tiendanube")
					return
				}
				logger.Info(bgCtx).Str("integration_id", integrationID).Msg("Webhooks creados automaticamente para Tiendanube")
			}()
		})
	} else {
		logger.Warn(context.Background()).Msg("Ni WEBHOOK_BASE_URL ni URL_BASE_SWAGGER configuradas, no se crearan webhooks automaticamente para Tiendanube")
	}

	return tncore.New(uc)
}
