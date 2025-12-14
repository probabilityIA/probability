package shopify

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/client"
	shopifycore "github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type integrationServiceAdapter struct {
	coreIntegration core.IIntegrationCore
}

func (a *integrationServiceAdapter) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	integration, err := a.coreIntegration.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	return &domain.Integration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integration.IntegrationType,
		Config:          integration.Config,
	}, nil
}

func (a *integrationServiceAdapter) GetIntegrationByStoreID(ctx context.Context, storeID string) (*domain.Integration, error) {
	integration, err := a.coreIntegration.GetIntegrationByStoreID(ctx, storeID, core.IntegrationTypeShopify)
	if err != nil {
		return nil, err
	}

	return &domain.Integration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integration.IntegrationType,
		Config:          integration.Config,
	}, nil
}

func (a *integrationServiceAdapter) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return a.coreIntegration.DecryptCredential(ctx, integrationID, fieldName)
}

func New(router *gin.RouterGroup, logger log.ILogger, config env.IConfig, coreIntegration core.IIntegrationCore, rabbitMQ rabbitmq.IQueue) {
	shopifyClient := client.New()

	// Crear publisher solo si RabbitMQ está disponible
	var orderPublisher domain.OrderPublisher
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
		logger.Info(context.Background()).
			Msg("RabbitMQ publisher initialized for Shopify integration")
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, Shopify orders will not be published to queue")
		// Crear un publisher no-op para evitar panics
		orderPublisher = queue.NewNoOpPublisher(logger)
	}

	shopifyCore := shopifycore.New(coreIntegration, shopifyClient, orderPublisher)
	coreIntegration.RegisterIntegration(core.IntegrationTypeShopify, shopifyCore)

	integrationService := &integrationServiceAdapter{
		coreIntegration: coreIntegration,
	}

	useCase := usecases.New(integrationService, shopifyClient, orderPublisher)

	shopifyHandler := handlers.New(useCase, logger, coreIntegration)
	shopifyHandler.RegisterRoutes(router, logger)
}
