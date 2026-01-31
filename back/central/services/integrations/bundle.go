package integrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/events"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify"
	whatsApp "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp"
	"github.com/secamc93/probability/back/central/services/modules"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
)

// IntegrationEventService es el servicio de eventos de integraciones (exportado para uso en otros módulos)
// Se accede a través de events.GetEventService()
var IntegrationEventService interface{}

// New inicializa todos los servicios de integraciones
// Este bundle coordina la inicialización de todos los módulos de integraciones
// (core, WhatsApp, Shopify, etc.) sin exponer dependencias externas
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, config env.IConfig, rabbitMQ rabbitmq.IQueue, s3 storage.IS3Service, redisClient redisclient.IRedis, moduleBundles *modules.ModuleBundles) {
	// Inicializar módulo de eventos de integraciones
	eventsRouter := router.Group("/integrations")
	eventService, _ := events.New(eventsRouter, logger)
	IntegrationEventService = eventService
	// Establecer instancia global para acceso desde otros módulos
	events.SetEventService(eventService)

	integrationCore := core.New(router, db, logger, config, s3)

	// Inicializar WhatsApp con configuración de notificaciones
	whatsappBundle := whatsApp.New(config, logger, db, rabbitMQ, redisClient, moduleBundles)

	integrationCore.RegisterIntegration(core.IntegrationTypeWhatsApp, whatsappBundle)

	shopify.New(router, logger, config, integrationCore, rabbitMQ, db)

}
