package woocommerce

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/primary/handlers"
	wooqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/client"
	woocore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/queue"
	wooproductrepo "github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/infra/secondary/repository"
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
	logger = logger.WithModule("woocommerce")

	httpClient := client.New()
	integrationService := woocore.NewIntegrationService(coreIntegration)
	productRepo := wooproductrepo.New(database, logger)

	var orderPublisher = queue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, WooCommerce orders will not be published to queue")
	}

	uc := usecases.New(httpClient, integrationService, orderPublisher, productRepo, rabbitMQ, logger)

	handler := handlers.New(uc, logger, rabbitMQ)
	handler.RegisterRoutes(router, logger)

	if rabbitMQ != nil {
		pushConsumer := wooqueue.NewInventoryPushConsumer(rabbitMQ, uc, logger)
		pushConsumer.Start(context.Background())

		productSyncConsumer := wooqueue.NewProductSyncConsumer(rabbitMQ, uc, logger)
		productSyncConsumer.Start(context.Background())

		webhookConsumer := wooqueue.NewWebhookConsumer(rabbitMQ, uc, logger)
		webhookConsumer.Start(context.Background())
	}

	baseURL := config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = config.Get("URL_BASE_SWAGGER")
	}
	if baseURL != "" {
		webhookSecret := config.Get("WOOCOMMERCE_WEBHOOK_SECRET")
		coreIntegration.OnIntegrationCreated(integrationcore.IntegrationTypeWoocommerce, func(obsCtx context.Context, integration *integrationcore.PublicIntegration) {
			go func() {
				bgCtx := context.Background()
				integrationID := fmt.Sprintf("%d", integration.ID)
				if err := uc.CreateWebhooks(bgCtx, integrationID, baseURL, webhookSecret); err != nil {
					logger.Error(bgCtx).Err(err).Str("integration_id", integrationID).Msg("Error al crear webhooks automaticamente para WooCommerce")
				} else {
					logger.Info(bgCtx).Str("integration_id", integrationID).Msg("Webhooks creados automaticamente para WooCommerce")
				}
			}()
		})
	} else {
		logger.Warn(context.Background()).Msg("Ni WEBHOOK_BASE_URL ni URL_BASE_SWAGGER configuradas, no se crearan webhooks automaticamente para WooCommerce")
	}

	return woocore.New(uc)
}
