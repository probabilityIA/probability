package tiktok

import (
	"context"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/primary/handlers"
	primaryqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/client"
	tiktokcore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/core"
	secondaryqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/repository"
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
	logger = logger.WithModule("tiktok")

	httpClient := client.New(logger)
	integrationService := tiktokcore.NewIntegrationService(coreIntegration)
	snapshotRepo := repository.New(database, logger)

	var publisher = secondaryqueue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		publisher = secondaryqueue.New(rabbitMQ, logger)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, TikTok snapshots will not be published to queue")
	}

	uc := usecases.New(httpClient, integrationService, publisher, snapshotRepo, logger)

	handler := handlers.New(uc, logger)
	handler.RegisterRoutes(router, logger)

	if rabbitMQ != nil {
		consumer := primaryqueue.NewSnapshotConsumer(rabbitMQ, uc, logger)
		consumer.Start(context.Background())
	}

	return tiktokcore.New(uc)
}
