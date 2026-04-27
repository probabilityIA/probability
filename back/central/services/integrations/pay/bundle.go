package pay

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold"
	"github.com/secamc93/probability/back/central/services/integrations/pay/epayco"
	"github.com/secamc93/probability/back/central/services/integrations/pay/melipago"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi"
	"github.com/secamc93/probability/back/central/services/integrations/pay/payu"
	"github.com/secamc93/probability/back/central/services/integrations/pay/router"
	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe"
	"github.com/secamc93/probability/back/central/services/integrations/pay/wompi"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	apiRouter *gin.RouterGroup,
	config env.IConfig,
	logger log.ILogger,
	database db.IDatabase,
	rabbitMQ rabbitmq.IQueue,
	coreSvc core.IIntegrationCore,
) {
	nequi.New(config, logger, database, rabbitMQ)
	bold.New(apiRouter, coreSvc, logger, rabbitMQ)
	wompi.New(config, logger, database, rabbitMQ)
	stripe.New(config, logger, database, rabbitMQ)
	payu.New(config, logger, database, rabbitMQ)
	epayco.New(config, logger, database, rabbitMQ)
	melipago.New(config, logger, database, rabbitMQ)

	router.New(logger, rabbitMQ)
}
