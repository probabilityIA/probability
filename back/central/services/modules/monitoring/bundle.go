package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa el módulo de monitoreo (relay puro, sin base de datos)
func New(router *gin.RouterGroup, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue) {
	logger = logger.WithModule("monitoring")

	// 1. Infraestructura secundaria (publisher de alertas a RabbitMQ)
	publisher := queue.New(rabbitMQ, logger)

	// 2. Capa de aplicación (caso de uso)
	useCase := app.New(publisher, logger, environment)

	// 3. Infraestructura primaria (handler HTTP)
	handler := handlers.New(useCase, logger, environment)

	// 4. Registrar rutas
	handler.RegisterRoutes(router)
}
