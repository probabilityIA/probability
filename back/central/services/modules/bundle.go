package modules

import (
	"github.com/gin-gonic/gin"
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
// NOTA: Los módulos deberían consultar datos directamente desde la BD o vía API HTTP,
// no a través de bundles. Este struct está vacío por diseño.
type ModuleBundles struct {
	// Vacío intencionalmente - los módulos no deben depender entre sí vía bundles
}

// New inicializa todos los módulos y retorna referencias a bundles compartidos
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
	notification_config.New(router, database, logger)

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

	// Inicializar módulo de invoicing (facturación electrónica)
	invoicing.New(router, database, logger, environment, rabbitMQ)

	// Retornar referencias a bundles compartidos
	return &ModuleBundles{
		// Vacío intencionalmente
	}
}
