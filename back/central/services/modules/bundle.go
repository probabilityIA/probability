package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus"
	"github.com/secamc93/probability/back/central/services/modules/payments"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	sharedQueue "github.com/secamc93/probability/back/central/shared/queue"
)

// New inicializa todos los módulos
func New(router *gin.RouterGroup, database db.IDatabase, logger log.ILogger, environment env.IConfig, rabbitMQ sharedQueue.IQueue) {
	// Inicializar módulo de payments
	payments.New(router, database, logger, environment)

	// Inicializar módulo de order status mappings
	orderstatus.New(router, database, logger, environment)

	// Inicializar módulo de orders
	orders.New(router, database, logger, environment, rabbitMQ)
}
