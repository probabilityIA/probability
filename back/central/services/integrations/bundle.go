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
	// 1. Inicializar Core (siempre necesario - registra rutas y expone interfaz pública)
	// El router ya viene como /api/v1, así que core.New registrará las rutas en /api/v1/integrations
	integrationCore := core.New(router, db, logger, config)

	whatsappBundle := whatsapp.New(config, logger)
	// Registrar con ambos códigos posibles para compatibilidad
	if err := integrationCore.RegisterTester(core.IntegrationTypeWhatsApp, whatsappBundle); err != nil {
		logger.Error().Err(err).Msg("Error registrando tester de WhatsApp")
	}
	// También registrar con "whatsap" (sin doble 'p') por si viene de la BD con ese código
	if err := integrationCore.RegisterTester("whatsap", whatsappBundle); err != nil {
		logger.Error().Err(err).Msg("Error registrando tester de WhatsApp (código alternativo)")
	}

	shopify.New(router, db, logger, config, integrationCore)
}
