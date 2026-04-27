package handlers

import (
	stderrors "errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	boldErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const boldSignatureHeader = "x-bold-signature"

type WebhookHandlers struct {
	useCase ports.IWebhookUseCase
	log     log.ILogger
}

func NewWebhookHandlers(useCase ports.IWebhookUseCase, logger log.ILogger) *WebhookHandlers {
	return &WebhookHandlers{
		useCase: useCase,
		log:     logger.WithModule("bold.webhook_handler"),
	}
}

func (h *WebhookHandlers) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/webhooks/bold", h.HandleWebhook)
}

func (h *WebhookHandlers) HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.log.Error(c.Request.Context()).Err(err).Msg("bold webhook read body failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	signature := c.GetHeader(boldSignatureHeader)

	if err := h.useCase.HandleIncomingWebhook(c.Request.Context(), signature, body); err != nil {
		switch {
		case stderrors.Is(err, boldErrors.ErrInvalidSignature):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature", "code": "BOLD_UNAUTHORIZED"})
		case stderrors.Is(err, boldErrors.ErrBoldConfigNotFound),
			stderrors.Is(err, boldErrors.ErrInvalidCredentials):
			h.log.Error(c.Request.Context()).Err(err).Msg("bold webhook: integration not configured")
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "bold not configured", "code": "BOLD_NOT_CONFIGURED"})
		default:
			h.log.Error(c.Request.Context()).Err(err).Msg("bold webhook processing failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal", "code": "INTERNAL"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
