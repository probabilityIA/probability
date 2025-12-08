package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/events"
	"github.com/secamc93/probability/back/central/services/modules/notification_config"
	"github.com/secamc93/probability/back/central/services/modules/orders"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus"
	"github.com/secamc93/probability/back/central/services/modules/payments"
	"github.com/secamc93/probability/back/central/services/modules/products"
	"github.com/secamc93/probability/back/central/services/modules/shipments"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// New inicializa todos los módulos
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ rabbitmq.IQueue, redisClient redis.IRedis) {
	// Inicializar módulo de payments
	payments.New(router, database, logger, environment)

	// Inicializar módulo de order status mappings
	orderstatus.New(router, database, logger, environment)

	// Inicializar módulo de orders
	orders.New(router, database, logger, environment, rabbitMQ, redisClient)

	// Inicializar módulo de products
	products.New(router, database, logger, environment)

	// Inicializar módulo de shipments
	shipments.New(router, database, logger, environment)

	// Inicializar módulo de notification configs
	notification_config.New(router, database)

	// Inicializar módulo de events (notificaciones en tiempo real)
	if redisClient != nil {
		events.New(router, database, logger, environment, redisClient)
	} else {
		logger.Warn().
			Msg("Redis no disponible, módulo de eventos no se inicializará")
	}
}
