package ecommerce

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/magento"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
	database db.IDatabase,
	integrationCore core.IIntegrationCore,
) {
	shopify.New(router, logger, config, integrationCore, rabbitMQ, database)

	meliProvider := meli.New(router, logger, config, rabbitMQ, database, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeMercadoLibre, meliProvider)

	wooProvider := woocommerce.New(router, logger, config, rabbitMQ, database, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeWoocommerce, wooProvider)

	vtexProvider := vtex.New(router, logger, config, rabbitMQ, database, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeVTEX, vtexProvider)

	tiendanubeProvider := tiendanube.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeTiendanube, tiendanubeProvider)

	magentoProvider := magento.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeMagento, magentoProvider)

	amazonProvider := amazon.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeAmazon, amazonProvider)

	falabellaProvider := falabella.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeFalabella, falabellaProvider)

	exitoProvider := exito.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeExito, exitoProvider)

	jumpsellerProvider := jumpseller.New(router, logger, config, rabbitMQ, database, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeJumpseller, jumpsellerProvider)
}
