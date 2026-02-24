package vtex

import (
	"context"

	"github.com/gin-gonic/gin"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/client"
	vtexcore "github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de VTEX y retorna el provider para registrar en integrationCore.
// type_id = 16 (IntegrationTypeVTEX)
func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
	coreIntegration integrationcore.IIntegrationCore,
) integrationcore.IIntegrationContract {
	logger = logger.WithModule("vtex")

	// 1. Infraestructura secundaria
	httpClient := client.New()
	integrationService := vtexcore.NewIntegrationService(coreIntegration)

	// Publisher de órdenes a RabbitMQ (con fallback no-op si no hay conexión)
	var orderPublisher = queue.NewNoOpPublisher(logger)
	if rabbitMQ != nil {
		orderPublisher = queue.New(rabbitMQ, logger, config)
	} else {
		logger.Warn(context.Background()).
			Msg("RabbitMQ not available, VTEX orders will not be published to queue")
	}

	// 2. Casos de uso
	uc := usecases.New(httpClient, integrationService, orderPublisher, logger)

	// 3. Handlers HTTP
	handler := handlers.New(uc, logger)
	handler.RegisterRoutes(router, logger)

	// 4. Retornar provider para que el bundle padre lo registre en el core
	return vtexcore.New(uc)
}
