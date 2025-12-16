package handlerintegrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// ListWebhooksHandler lista todos los webhooks de una integración
//
//	@Summary		Listar webhooks
//	@Description	Lista todos los webhooks configurados para una integración (solo disponible para integraciones que lo soporten, como Shopify)
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID de la integración"
//	@Success		200	{object}	response.ListWebhooksResponse
//	@Failure		400	{object}	response.ErrorResponse	"ID inválido"
//	@Failure		404	{object}	response.ErrorResponse	"Integración no encontrada"
//	@Failure		500	{object}	response.ErrorResponse	"Error interno"
//	@Router			/integrations/{id}/webhooks [get]
func (h *IntegrationHandler) ListWebhooksHandler(c *gin.Context) {
	idStr := c.Param("id")

	// Listar webhooks a través del core
	webhooks, err := h.orderSyncSvc.ListWebhooks(c.Request.Context(), idStr)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Error al listar webhooks")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ListWebhooksResponse{
		Success: true,
		Data:    webhooks,
	})
}

