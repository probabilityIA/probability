package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/transport/enviame"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick"
	"github.com/secamc93/probability/back/central/services/integrations/transport/mipaquete"
	"github.com/secamc93/probability/back/central/services/integrations/transport/router"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	ginRouter *gin.RouterGroup,
	database db.IDatabase,
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	integrationCore core.IIntegrationCore,
) {
	envioclick.New(ginRouter, database, logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeEnvioClick, nil)

	enviame.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeEnviame, nil)

	mipaquete.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeMiPaquete, nil)

	router.New(logger, rabbitMQ)
}
