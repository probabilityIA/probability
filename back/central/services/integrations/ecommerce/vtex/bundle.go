package vtex

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/primary/handlers"
	vtexqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/primary/queue"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/client"
	vtexcore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/queue"
	vtexproductrepo "github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/repository"
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
	logger = logger.WithModule("vtex")

	httpClient := client.New()
	integrationService := vtexcore.NewIntegrationService(coreIntegration)
	productRepo := vtexproductrepo.New(database, logger)

	var orderPublisher = queue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, VTEX orders will not be published to queue")
	}

	baseURL := config.Get("WEBHOOK_BASE_URL")
	if baseURL == "" {
		baseURL = config.Get("URL_BASE_SWAGGER")
	}

	uc := usecases.New(httpClient, integrationService, orderPublisher, productRepo, rabbitMQ, baseURL, logger)

	handler := handlers.New(uc, baseURL, logger)
	handler.RegisterRoutes(router, logger)

	pushConsumer := vtexqueue.NewInventoryPushConsumer(rabbitMQ, uc, logger)
	pushConsumer.Start(context.Background())

	if baseURL != "" {
		coreIntegration.OnIntegrationCreated(integrationcore.IntegrationTypeVTEX, func(obsCtx context.Context, integration *integrationcore.PublicIntegration) {
			go func() {
				bgCtx := context.Background()
				integrationID := strconv.FormatUint(uint64(integration.ID), 10)
				if err := uc.CreateWebhooks(bgCtx, integrationID, baseURL, false); err != nil {
					if errors.Is(err, domain.ErrForeignHookExists) {
						logger.Warn(bgCtx).Err(err).
							Uint("integration_id", integration.ID).
							Msg("La cuenta VTEX ya tiene un hook de otra herramienta: registralo manualmente desde la integracion si quieres reemplazarlo")
						return
					}
					logger.Error(bgCtx).Err(err).
						Uint("integration_id", integration.ID).
						Msg("Error al registrar el hook de ordenes en VTEX")
				}
			}()
		})
	} else {
		logger.Warn(context.Background()).
			Msg("WEBHOOK_BASE_URL no configurada, VTEX no registrara hooks automaticamente")
	}

	return vtexcore.New(uc)
}
