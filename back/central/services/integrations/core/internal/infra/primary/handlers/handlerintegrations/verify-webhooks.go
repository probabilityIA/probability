package handlerintegrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// VerifyWebhooksHandler verifica webhooks existentes que coincidan con nuestra URL
//
//	@Summary		Verificar webhooks por URL
//	@Description	Verifica si existen webhooks en la plataforma externa que coincidan con nuestra URL generada. Solo retorna webhooks que coinciden exactamente con nuestra URL, sin afectar otros webhooks del cliente.
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID de la integración"
//	@Success		200	{object}	response.VerifyWebhooksResponse
//	@Failure		400	{object}	response.ErrorResponse	"ID inválido"
//	@Failure		404	{object}	response.ErrorResponse	"Integración no encontrada"
//	@Failure		500	{object}	response.ErrorResponse	"Error interno"
//	@Router			/integrations/{id}/webhooks/verify [get]
func (h *IntegrationHandler) VerifyWebhooksHandler(c *gin.Context) {
	idStr := c.Param("id")

	// Verificar webhooks a través del core
	webhooks, err := h.orderSyncSvc.VerifyWebhooksByURL(c.Request.Context(), idStr)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Error al verificar webhooks")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.VerifyWebhooksResponse{
		Success: true,
		Data:    webhooks,
		Message: "Webhooks verificados exitosamente",
	})
}
