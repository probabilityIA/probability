package handlerintegrations

import (
	"encoding/json"
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
	webhooks, err := h.usecase.ListWebhooks(c.Request.Context(), idStr)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Error al listar webhooks")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ListWebhooksResponse{
		Success: true,
		Data:    mapWebhooksToResponse(webhooks),
	})
}

// mapWebhooksToResponse convierte []interface{} (domain structs sin JSON tags) a []WebhookInfoResponse (con JSON tags snake_case)
func mapWebhooksToResponse(webhooks []interface{}) []response.WebhookInfoResponse {
	if len(webhooks) == 0 {
		return []response.WebhookInfoResponse{}
	}

	// Los domain structs no tienen JSON tags, se serializan como PascalCase
	data, err := json.Marshal(webhooks)
	if err != nil {
		return []response.WebhookInfoResponse{}
	}

	var rawList []map[string]interface{}
	if err := json.Unmarshal(data, &rawList); err != nil {
		return []response.WebhookInfoResponse{}
	}

	result := make([]response.WebhookInfoResponse, 0, len(rawList))
	for _, raw := range rawList {
		result = append(result, response.WebhookInfoResponse{
			ID:        getStringField(raw, "ID"),
			Address:   getStringField(raw, "Address"),
			Topic:     getStringField(raw, "Topic"),
			Format:    getStringField(raw, "Format"),
			CreatedAt: getStringField(raw, "CreatedAt"),
			UpdatedAt: getStringField(raw, "UpdatedAt"),
		})
	}
	return result
}

func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
