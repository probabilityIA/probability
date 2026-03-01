package integrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing"
	"github.com/secamc93/probability/back/central/services/integrations/messaging"
	pay "github.com/secamc93/probability/back/central/services/integrations/pay"
	"github.com/secamc93/probability/back/central/services/integrations/transport"
	"github.com/secamc93/probability/back/central/services/modules"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// New inicializa todos los servicios de integraciones.
// Retorna core.IIntegrationCore para que otros módulos puedan usarlo.
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, config env.IConfig, rabbitMQ rabbitmq.IQueue, s3 storage.IS3Service, redisClient redisclient.IRedis, moduleBundles *modules.ModuleBundles) core.IIntegrationCore {
	// Events publisher se inicializa en init.go (módulo unificado services/events)

	// Inicializar Integration Core (hub central de integraciones)
	integrationCore := core.New(router, db, redisClient, logger, config, s3)

	// ═══════════════════════════════════════════════════════════════
	// REGISTRO DE INTEGRACIONES
	// ═══════════════════════════════════════════════════════════════

	// Messaging: todos los proveedores de mensajería
	messaging.New(config, logger, db, rabbitMQ, redisClient, moduleBundles, integrationCore)

	// E-commerce: todos los proveedores de e-commerce
	ecommerce.New(router, logger, config, rabbitMQ, db, integrationCore)

	// Invoicing: todos los proveedores de facturación electrónica + router de colas
	invoicing.New(config, logger, rabbitMQ, integrationCore)

	// Transport: todos los proveedores de transporte + router de colas
	transport.New(logger, rabbitMQ, integrationCore)

	// Pay: todos los proveedores de pago (Nequi, etc.) + router de colas
	pay.New(config, logger, db, rabbitMQ)

	return integrationCore
}
