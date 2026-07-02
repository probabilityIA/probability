package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type Handler struct {
	webhookLog      ports.IWebhookLogRepository
	coreIntegration core.IIntegrationService
	rabbit          rabbitmq.IQueue
	log             log.ILogger
}

func New(
	webhookLog ports.IWebhookLogRepository,
	coreIntegration core.IIntegrationService,
	rabbit rabbitmq.IQueue,
	logger log.ILogger,
) *Handler {
	return &Handler{
		webhookLog:      webhookLog,
		coreIntegration: coreIntegration,
		rabbit:          rabbit,
		log:             logger.WithModule("siigo.webhook_handler"),
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/siigo")
	{
		group.POST("/webhook", h.HandleWebhook)
	}
}
