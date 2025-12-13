package integrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify"
	ordersdomain "github.com/secamc93/probability/back/central/services/modules/orders/domain"

	// "github.com/secamc93/probability/back/central/services/integrations/test"
	whatsapp "github.com/secamc93/probability/back/central/services/integrations/whatsApp"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa todos los servicios de integraciones
// Este bundle coordina la inicialización de todos los módulos de integraciones
// (core, WhatsApp, Shopify, etc.) sin exponer dependencias externas
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, config env.IConfig, rabbitMQ rabbitmq.IQueue, orderUseCase ordersdomain.IOrderMappingUseCase) {

	integrationCore := core.New(router, db, logger, config)

	whatsappBundle := whatsapp.New(config, logger)

	integrationCore.RegisterTester(core.IntegrationTypeWhatsApp, whatsappBundle)

	integrationCore.RegisterTester("whatsap", whatsappBundle)

	shopify.New(router, db, logger, config, integrationCore, orderUseCase)

	// test.New(router, logger, rabbitMQ)
}
