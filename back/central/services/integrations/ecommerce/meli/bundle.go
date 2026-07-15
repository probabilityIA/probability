package meli

import (
	"context"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/primary/handlers"
	meliqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/client"
	melicore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/infra/secondary/repository"
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
	logger = logger.WithModule("meli")

	httpClient := client.New()
	integrationService := melicore.NewIntegrationService(coreIntegration)
	productRepo := repository.New(database, logger)
	inventoryRepo := repository.NewInventory(database, logger)
	orderLookupRepo := repository.NewOrderLookup(database, logger)

	var orderPublisher = queue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, MercadoLibre orders will not be published to queue")
	}

	uc := usecases.New(httpClient, integrationService, orderPublisher, productRepo, inventoryRepo, rabbitMQ, logger)

	handler := handlers.New(uc, logger, config, coreIntegration)
	handler.RegisterRoutes(router, logger)

	if rabbitMQ != nil {
		statusConsumer := meliqueue.NewOrderStatusConsumer(rabbitMQ, uc, orderLookupRepo, logger)
		statusConsumer.Start(context.Background())

		billingConsumer := meliqueue.NewBillingRetryConsumer(rabbitMQ, uc, logger)
		billingConsumer.Start(context.Background())
	}

	return melicore.New(uc)
}
