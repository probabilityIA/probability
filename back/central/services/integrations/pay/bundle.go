package pay

import (
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi"
	"github.com/secamc93/probability/back/central/services/integrations/pay/router"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa todas las integraciones de pago y el router
func New(
	config env.IConfig,
	logger log.ILogger,
	database db.IDatabase,
	rabbitMQ rabbitmq.IQueue,
) {
	// Inicializar Nequi (gateway de pago)
	nequi.New(config, logger, database, rabbitMQ)

	// Router: consume pay.requests y enruta al gateway correcto
	// Se inicializa al final para que las colas de gateways ya estén declaradas
	router.New(logger, rabbitMQ)
}
