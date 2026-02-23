package integrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce"
	"github.com/secamc93/probability/back/central/services/integrations/events"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing"
	"github.com/secamc93/probability/back/central/services/integrations/messaging"
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
	// Inicializar publisher de eventos de integraciones (publica a Redis)
	// La entrega SSE al frontend la maneja modules/events (centralizada)
	events.Init(logger, redisClient)

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

	return integrationCore
}
