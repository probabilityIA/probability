package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define todos los métodos HTTP de los handlers de WhatsApp
type IHandler interface {
	// Rutas
	RegisterRoutes(router *gin.RouterGroup)

	// Template
	SendTemplate(c *gin.Context)

	// Webhook
	VerifyWebhook(c *gin.Context)
	ReceiveWebhook(c *gin.Context)
}

// handler contiene las dependencias compartidas
type Handler struct {
	useCase usecasemessaging.IUseCase
	log     log.ILogger
	config  env.IConfig
}

// New crea la instancia única de handler con todas las dependencias
func New(
	useCase usecasemessaging.IUseCase,
	logger log.ILogger,
	config env.IConfig,
) IHandler {
	return &Handler{
		useCase: useCase,
		log:     logger,
		config:  config,
	}
}
