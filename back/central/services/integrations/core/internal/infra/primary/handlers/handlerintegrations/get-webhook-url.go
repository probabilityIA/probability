package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// GetWebhookURLHandler obtiene la URL del webhook para una integración
//
//	@Summary		Obtener URL del webhook
//	@Description	Obtiene la URL del webhook que debe configurarse en la plataforma externa (ej: Shopify) para recibir eventos en tiempo real
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID de la integración"
//	@Success		200	{object}	response.WebhookURLResponse
//	@Failure		400	{object}	response.ErrorResponse	"ID inválido"
//	@Failure		404	{object}	response.ErrorResponse	"Integración no encontrada"
//	@Failure		500	{object}	response.ErrorResponse	"Error interno"
//	@Router			/integrations/{id}/webhook [get]
func (h *IntegrationHandler) GetWebhookURLHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("ID de integración inválido")
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "ID de integración inválido",
		})
		return
	}

	// Obtener la información del webhook a través del servicio de sincronización
	// que tiene acceso al core de integraciones
	webhookInfo, err := h.usecase.GetWebhookURL(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("Error al obtener URL del webhook")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebhookURLResponse{
		Success: true,
		Data: &response.WebhookURLData{
			URL:         webhookInfo.URL,
			Method:      webhookInfo.Method,
			Description: webhookInfo.Description,
			Events:      webhookInfo.Events,
		},
	})
}
