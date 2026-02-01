package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/ai"
	"github.com/secamc93/probability/back/central/services/modules/dashboard"
	"github.com/secamc93/probability/back/central/services/modules/events"
	"github.com/secamc93/probability/back/central/services/modules/fulfillmentstatus"
	"github.com/secamc93/probability/back/central/services/modules/invoicing"
	"github.com/secamc93/probability/back/central/services/modules/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/orders"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus"
	"github.com/secamc93/probability/back/central/services/modules/payments"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus"
	"github.com/secamc93/probability/back/central/services/modules/products"
	"github.com/secamc93/probability/back/central/services/modules/shipments"
	"github.com/secamc93/probability/back/central/services/modules/wallet"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// ModuleBundles contiene referencias a los bundles de módulos que otros servicios pueden necesitar
type ModuleBundles struct {
	router          *gin.RouterGroup
	database        db.IDatabase
	logger          log.ILogger
	environment     env.IConfig
	rabbitMQ        rabbitmq.IQueue
	integrationCore core.IIntegrationCore
}

// New inicializa todos los módulos (excepto invoicing que requiere integrationCore)
// y retorna referencias a bundles compartidos
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redis.IRedis) *ModuleBundles {
	// Inicializar módulo de payments
	payments.New(router, database, logger, environment)

	// Inicializar módulo de order status mappings
	orderstatus.New(router, database, logger, environment)

	// Inicializar módulo de payment statuses
	paymentstatus.New(router, database, logger, environment)

	// Inicializar módulo de fulfillment statuses
	fulfillmentstatus.New(router, database, logger, environment)

	// Inicializar módulo de orders
	orders.New(router, database, logger, environment, rabbitMQ, redisClient)

	// Inicializar módulo de products
	products.New(router, database, logger, environment)

	// Inicializar módulo de shipments
	shipments.New(router, database, logger, environment)

	// Inicializar módulo de notification configs
	notification_config.New(router, database, redisClient, logger)

	// Inicializar módulo de events (notificaciones en tiempo real)
	if redisClient != nil {
		events.New(router, database, logger, environment, redisClient)
	} else {
		logger.Warn().
			Msg("Redis no disponible, módulo de eventos no se inicializará")
	}
	// Inicializar módulo de AI
	ai.New(router, logger)

	// Inicializar módulo de dashboard
	dashboard.New(router, database, logger)

	// Inicializar módulo de wallet
	wallet.New(router, database, logger, environment)

	// NOTA: invoicing se inicializa en SetIntegrationCore() porque depende de integrationCore

	// Retornar referencias a bundles compartidos
	return &ModuleBundles{
		router:      router,
		database:    database,
		logger:      logger,
		environment: environment,
		rabbitMQ:    rabbitMQ,
	}
}

// SetIntegrationCore establece el integrationCore y luego inicializa módulos que lo requieren
func (mb *ModuleBundles) SetIntegrationCore(integrationCore core.IIntegrationCore) {
	mb.integrationCore = integrationCore

	// Inicializar módulo de invoicing (requiere integrationCore)
	invoicing.New(mb.router, mb.database, mb.logger, mb.environment, mb.rabbitMQ, integrationCore)
}
