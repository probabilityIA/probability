package tiendanube

import (
	"context"

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

	return tncore.New(uc)
}
