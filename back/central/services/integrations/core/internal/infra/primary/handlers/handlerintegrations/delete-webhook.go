package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// DeleteWebhookHandler elimina un webhook de una integración
//
//	@Summary		Eliminar webhook
//	@Description	Elimina un webhook específico de una integración (solo disponible para integraciones que lo soporten, como Shopify)
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		int		true	"ID de la integración"
//	@Param			webhook_id	query		string	true	"ID del webhook a eliminar"
//	@Success		200			{object}	response.DeleteWebhookResponse
//	@Failure		400			{object}	response.ErrorResponse	"ID inválido o webhook_id faltante"
//	@Failure		404			{object}	response.ErrorResponse	"Integración o webhook no encontrado"
//	@Failure		500			{object}	response.ErrorResponse	"Error interno"
//	@Router			/integrations/{id}/webhooks/{webhook_id} [delete]
func (h *IntegrationHandler) DeleteWebhookHandler(c *gin.Context) {
	idStr := c.Param("id")
	webhookID := c.Param("webhook_id")

	if webhookID == "" {
		h.logger.Error().Str("id", idStr).Msg("webhook_id es requerido")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "webhook_id es requerido",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("ID de integración inválido")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "ID de integración inválido",
		})
		return
	}

	_ = id // Evitar warning de variable no usada

	// Eliminar webhook a través del core
	if err := h.usecase.DeleteWebhook(c.Request.Context(), idStr, webhookID); err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Str("webhook_id", webhookID).Msg("Error al eliminar webhook")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.DeleteWebhookResponse{
		Success: true,
		Message: "Webhook eliminado exitosamente",
	})
}













