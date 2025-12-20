package handlerintegrations

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// CreateWebhookHandler crea webhooks en la plataforma externa
//
//	@Summary		Crear webhooks
//	@Description	Crea webhooks en la plataforma externa (ej: Shopify) después de verificar y eliminar webhooks duplicados que coincidan con nuestra URL. Solo elimina webhooks que coinciden exactamente con nuestra URL, sin afectar otros webhooks del cliente.
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID de la integración"
//	@Success		200	{object}	response.CreateWebhookResponse
//	@Failure		400	{object}	response.ErrorResponse	"ID inválido"
//	@Failure		404	{object}	response.ErrorResponse	"Integración no encontrada"
//	@Failure		500	{object}	response.ErrorResponse	"Error interno"
//	@Router			/integrations/{id}/webhooks/create [post]
func (h *IntegrationHandler) CreateWebhookHandler(c *gin.Context) {
	idStr := c.Param("id")

	// Crear webhooks a través del core
	result, err := h.orderSyncSvc.CreateWebhook(c.Request.Context(), idStr)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Error al crear webhooks")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Convertir el resultado del dominio (sin etiquetas JSON) a la estructura de respuesta (con etiquetas JSON)
	// El resultado viene como interface{} desde el core, necesitamos convertirlo correctamente
	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Error al serializar resultado de creación de webhooks")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Error al procesar resultado de creación de webhooks",
		})
		return
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(resultJSON, &resultMap); err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("Error al deserializar resultado de creación de webhooks")
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Error al procesar resultado de creación de webhooks",
		})
		return
	}

	// Mapear desde nombres PascalCase (dominio sin etiquetas) a la estructura de respuesta
	data := response.CreateWebhookResponseData{
		ExistingWebhooks: []interface{}{},
		DeletedWebhooks:  []interface{}{},
		CreatedWebhooks:  []string{},
		WebhookURL:       "",
	}

	// Sin etiquetas JSON en el dominio, los campos se serializan con nombres PascalCase
	if existing, ok := resultMap["ExistingWebhooks"].([]interface{}); ok {
		data.ExistingWebhooks = existing
	}
	if deleted, ok := resultMap["DeletedWebhooks"].([]interface{}); ok {
		data.DeletedWebhooks = deleted
	}
	if created, ok := resultMap["CreatedWebhooks"].([]interface{}); ok {
		createdStrings := make([]string, len(created))
		for i, v := range created {
			if str, ok := v.(string); ok {
				createdStrings[i] = str
			}
		}
		data.CreatedWebhooks = createdStrings
	}
	if url, ok := resultMap["WebhookURL"].(string); ok {
		data.WebhookURL = url
	}

	c.JSON(http.StatusOK, response.CreateWebhookResponse{
		Success: true,
		Data:    data,
		Message: "Webhooks creados exitosamente",
	})
}
