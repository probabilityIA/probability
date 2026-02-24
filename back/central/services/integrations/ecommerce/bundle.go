package ecommerce

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella"
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

// New inicializa todos los proveedores de e-commerce y los registra en integrationCore.
// Debe llamarse después de inicializar integrationCore.
func New(
	router *gin.RouterGroup,
	logger log.ILogger,
	config env.IConfig,
	rabbitMQ rabbitmq.IQueue,
	database db.IDatabase,
	integrationCore core.IIntegrationCore,
) {
	// Shopify (type_id=1) — se auto-registra internamente (incluye OnIntegrationCreated)
	shopify.New(router, logger, config, integrationCore, rabbitMQ, database)

	// MercadoLibre (type_id=3)
	meliProvider := meli.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeMercadoLibre, meliProvider)

	// WooCommerce (type_id=4)
	wooProvider := woocommerce.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeWoocommerce, wooProvider)

	// VTEX (type_id=16)
	vtexProvider := vtex.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeVTEX, vtexProvider)

	// Tiendanube (type_id=17)
	tiendanubeProvider := tiendanube.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeTiendanube, tiendanubeProvider)

	// Magento (type_id=18)
	magentoProvider := magento.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeMagento, magentoProvider)

	// Amazon (type_id=19)
	amazonProvider := amazon.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeAmazon, amazonProvider)

	// Falabella (type_id=20)
	falabellaProvider := falabella.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeFalabella, falabellaProvider)

	// Éxito (type_id=21)
	exitoProvider := exito.New(router, logger, config, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeExito, exitoProvider)
}
