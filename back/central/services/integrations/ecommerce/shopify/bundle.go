package shopify

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/client"
	shopifycore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/eventpublisher"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(router *gin.RouterGroup, logger log.ILogger, config env.IConfig, coreIntegration core.IIntegrationCore, rabbitMQ rabbitmq.IQueue, database db.IDatabase) {
	shopifyClient := client.New()

	// Habilitar debug del cliente HTTP si está configurado
	debugMode := os.Getenv("SHOPIFY_DEBUG")
	if debugMode == "true" || debugMode == "1" {
		shopifyClient.SetDebug(true)
		logger.Info(context.Background()).Msg("🔍 Shopify HTTP client debug mode ENABLED")
	}

	// Crear publisher solo si RabbitMQ está disponible
	var orderPublisher domain.OrderPublisher
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, Shopify orders will not be published to queue")
		// Crear un publisher no-op para evitar panics
		orderPublisher = queue.NewNoOpPublisher(logger)
	}

	syncEventPub := eventpublisher.New(rabbitMQ)
	useCase := usecases.New(coreIntegration, shopifyClient, orderPublisher, logger, syncEventPub)

	shopifyCore := shopifycore.New(useCase)
	coreIntegration.RegisterIntegration(core.IntegrationTypeShopify, shopifyCore)

	// Registrar observador para crear webhook automáticamente cuando se crea una integración de Shopify
	baseURL := config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = config.Get("URL_BASE_SWAGGER")
	}

	if baseURL != "" {
		coreIntegration.OnIntegrationCreated(core.IntegrationTypeShopify, func(obsCtx context.Context, integration *core.PublicIntegration) {
			// Crear webhook de forma asíncrona para no bloquear la respuesta
			go func() {
				bgCtx := context.Background()
				integrationID := fmt.Sprintf("%d", integration.ID)
				_, err := useCase.CreateWebhook(bgCtx, integrationID, baseURL)
				if err != nil {
					logger.Error(bgCtx).
						Err(err).
						Str("integration_id", integrationID).
						Msg("Error al crear webhook automáticamente para integración de Shopify")
				} else {
					logger.Info(bgCtx).
						Str("integration_id", integrationID).
						Msg("Webhook creado automáticamente para integración de Shopify")
				}
			}()
		})
	} else {
		logger.Warn(context.Background()).
			Msg("Ni WEBHOOK_BASE_URL ni URL_BASE_SWAGGER están configuradas, no se crearán webhooks automáticamente para Shopify")
	}

	shopifyHandler := handlers.New(useCase, logger, config, coreIntegration)
	shopifyHandler.RegisterRoutes(router, logger)
}
