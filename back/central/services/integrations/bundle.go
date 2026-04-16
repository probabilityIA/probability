package integrations

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing"
	"github.com/secamc93/probability/back/central/services/integrations/messaging"
	pay "github.com/secamc93/probability/back/central/services/integrations/pay"
	storefrontprovider "github.com/secamc93/probability/back/central/services/integrations/storefront"
	websiteprovider "github.com/secamc93/probability/back/central/services/integrations/website"
	"github.com/secamc93/probability/back/central/services/integrations/transport"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/email"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
	"github.com/secamc93/probability/back/central/shared/storage"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// New inicializa todos los servicios de integraciones.
// Retorna core.IIntegrationCore para que otros módulos puedan usarlo.
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, config env.IConfig, rabbitMQ rabbitmq.IQueue, s3 storage.IS3Service, redisClient redisclient.IRedis, emailService email.IEmailService) core.IIntegrationCore {
	// Events publisher se inicializa en init.go (módulo unificado services/events)

	// Inicializar Integration Core (hub central de integraciones)
	integrationCore := core.New(router, db, redisClient, logger, config, s3, rabbitMQ)

	// REGISTRO DE INTEGRACIONES

	// Messaging: todos los proveedores de mensajería (sin DB — cache-first, DB-async)
	messaging.New(config, logger, rabbitMQ, redisClient, integrationCore, emailService, router)

	// E-commerce: todos los proveedores de e-commerce
	ecommerce.New(router, logger, config, rabbitMQ, db, integrationCore)

	// Invoicing: todos los proveedores de facturación electrónica + router de colas
	invoicing.New(config, logger, rabbitMQ, integrationCore)

	// Transport: todos los proveedores de transporte + router de colas
	transport.New(logger, rabbitMQ, integrationCore)

	// Pay: todos los proveedores de pago (Nequi, etc.) + router de colas
	pay.New(config, logger, db, rabbitMQ)

	// Storefront: Tienda y Tienda Web
	integrationCore.RegisterIntegration(core.IntegrationTypeTienda, storefrontprovider.New())
	integrationCore.RegisterIntegration(core.IntegrationTypeTiendaWeb, websiteprovider.New())

	// Auto-crear BusinessWebsiteConfig al crear integración Tienda Web
	integrationCore.OnIntegrationCreated(core.IntegrationTypeTiendaWeb, func(ctx context.Context, integration *core.PublicIntegration) {
		if integration.BusinessID == nil {
			return
		}
		var existing models.BusinessWebsiteConfig
		if err := db.Conn(ctx).Where("business_id = ?", *integration.BusinessID).First(&existing).Error; err == gorm.ErrRecordNotFound {
			websiteConfig := models.BusinessWebsiteConfig{
				BusinessID:           *integration.BusinessID,
				ShowHero:             true,
				ShowFeaturedProducts: true,
				ShowFullCatalog:      true,
				ShowContact:          true,
			}
			db.Conn(ctx).Create(&websiteConfig)
		}
	})

	return integrationCore
}
