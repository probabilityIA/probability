package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// IHandler define todos los métodos HTTP de los handlers de WhatsApp
type IHandler interface {
	// Rutas
	RegisterRoutes(router *gin.RouterGroup)

	// Template
	SendTemplate(c *gin.Context)

	// Manual reply desde dashboard
	SendManualReply(c *gin.Context)

	// AI control
	PauseAI(c *gin.Context)
	ResumeAI(c *gin.Context)

	// Webhook
	VerifyWebhook(c *gin.Context)
	ReceiveWebhook(c *gin.Context)

	// Post-init setter
	SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter)
}

// handler contiene las dependencias compartidas
type handler struct {
	useCase         usecasemessaging.IUseCase
	log             log.ILogger
	config          env.IConfig
	platformCredsGetter ports.IPlatformCredentialsGetter
}

// New crea la instancia única de handler con todas las dependencias
func New(
	useCase usecasemessaging.IUseCase,
	logger log.ILogger,
	config env.IConfig,
) IHandler {
	return &handler{
		useCase: useCase,
		log:     logger,
		config:  config,
	}
}

// SetPlatformCredsGetter inyecta el getter de credenciales de plataforma (post-init)
func (h *handler) SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter) {
	h.platformCredsGetter = getter
}
