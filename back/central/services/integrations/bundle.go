package integrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify"
	whatsapp "github.com/secamc93/probability/back/central/services/integrations/whatsApp"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// New inicializa todos los servicios de integraciones
// Este bundle coordina la inicialización de todos los módulos de integraciones
// (core, WhatsApp, Shopify, etc.) sin exponer dependencias externas
func New(router *gin.RouterGroup, db db.IDatabase, logger log.ILogger, config env.IConfig) {
	// Crear grupo /integrations para todas las rutas de integraciones
	integrationsGroup := router.Group("/integrations")

	// 1. Inicializar Core (siempre necesario - registra rutas y expone interfaz pública)
	integrationCore := core.New(integrationsGroup, db, logger, config)

	// 2. Inicializar módulos de integraciones específicos
	// Cada módulo se registra automáticamente con el core si es necesario
	_ = whatsapp.New(config, integrationCore)

	// Inicializar Shopify
	shopify.New(integrationsGroup, db, logger, config, integrationCore)
}
