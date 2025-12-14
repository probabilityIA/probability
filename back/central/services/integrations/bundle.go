package integrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa todos los servicios de integraciones
// Este bundle coordina la inicialización de todos los módulos de integraciones
// (core, WhatsApp, Shopify, etc.) sin exponer dependencias externas
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, config env.IConfig, rabbitMQ rabbitmq.IQueue) {

	integrationCore := core.New(router, db, logger, config)

	whatsappBundle := whatsApp.New(config, logger)

	integrationCore.RegisterIntegration(core.IntegrationTypeWhatsApp, whatsappBundle)

	shopify.New(router, logger, config, integrationCore, rabbitMQ)

}
