package handlerintegrations

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

//	@Summary		Listar webhooks
//	@Description	Lista todos los webhooks configurados para una integracion
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID de la integracion"
//	@Success		200	{object}	response.ListWebhooksResponse
//	@Failure		400	{object}	response.ErrorResponse	"ID invalido"
//	@Failure		404	{object}	response.ErrorResponse	"Integracion no encontrada"
//	@Failure		500	{object}	response.ErrorResponse	"Error interno"
//	@Router			/integrations/{id}/webhooks [get]
func (h *IntegrationHandler) ListWebhooksHandler(c *gin.Context) {
	idStr := c.Param("id")

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

func mapWebhooksToResponse(webhooks []interface{}) []response.WebhookInfoResponse {
	if len(webhooks) == 0 {
		return []response.WebhookInfoResponse{}
	}

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
		result = append(result, mapRawWebhookToResponse(raw))
	}
	return result
}

func mapRawWebhookToResponse(raw map[string]interface{}) response.WebhookInfoResponse {
	address := getStringField(raw, "Address", "address", "URL", "url")

	return response.WebhookInfoResponse{
		ID:        getStringField(raw, "ID", "id"),
		Address:   address,
		URL:       address,
		Topic:     getStringField(raw, "Topic", "topic"),
		Format:    getStringField(raw, "Format", "format"),
		Active:    getBoolField(raw, true, "Active", "active"),
		CreatedAt: getStringField(raw, "CreatedAt", "created_at"),
		UpdatedAt: getStringField(raw, "UpdatedAt", "updated_at"),
	}
}

func getStringField(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

func getBoolField(m map[string]interface{}, fallback bool, keys ...string) bool {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if b, ok := v.(bool); ok {
				return b
			}
		}
	}
	return fallback
}
