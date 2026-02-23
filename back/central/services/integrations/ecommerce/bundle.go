package ecommerce

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify"
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
}
