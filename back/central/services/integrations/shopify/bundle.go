package shopify

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/publisher"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/tester"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	sharedQueue "github.com/secamc93/probability/back/central/shared/queue"
)

func New(
	router *gin.RouterGroup,
	db db.IDatabase,
	logger log.ILogger,
	config env.IConfig,
	coreIntegration core.IIntegrationCore,
) {
	// 1. Init Secondary Adapters
	shopifyClient := client.New()

	// Init RabbitMQ connection
	rabbitMQ, err := sharedQueue.New(logger, config)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to connect to RabbitMQ, using log publisher as fallback")
		// Fallback to log publisher if RabbitMQ fails
		orderPublisher := publisher.New(logger)
		initializeModule(router, shopifyClient, orderPublisher, coreIntegration, logger)
		return
	}

	// Use RabbitMQ publisher
	orderPublisher := queue.New(rabbitMQ, logger)

	initializeModule(router, shopifyClient, orderPublisher, coreIntegration, logger)
}

func initializeModule(
	router *gin.RouterGroup,
	shopifyClient domain.ShopifyClient,
	orderPublisher domain.OrderPublisher,
	coreIntegration core.IIntegrationCore,
	logger log.ILogger,
) {
	// 2. Register Tester with Core
	shopifyTester := tester.New(shopifyClient)
	if err := coreIntegration.RegisterTester("shopify", shopifyTester); err != nil {
		logger.Error().Msg("Failed to register shopify tester: " + err.Error())
	}

	// 3. Init Use Cases
	syncUseCase := usecases.New(coreIntegration, shopifyClient, orderPublisher)

	// 4. Init Handlers
	h := handlers.New(syncUseCase)

	// 5. Register Routes
	h.RegisterRoutes(router)
}
