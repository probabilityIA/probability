package handlerintegrations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// GetIntegrationByTypeHandler obtiene una integración por tipo (para uso interno, como WhatsApp)
//
//	@Summary		Obtener integración por tipo
//	@Description	Obtiene una integración activa por tipo y business_id (para uso interno)
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			type		path		string	true	"Tipo de integración"	example("whatsapp")
//	@Param			business_id	query		int		false	"ID del business (opcional)"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
func (h *IntegrationHandler) GetIntegrationByTypeHandler(c *gin.Context) {
	integrationType := c.Param("type")
	if integrationType == "" {
		h.logger.Error().Str("endpoint", "/integrations/type/:type").Str("method", "GET").Msg("Tipo de integración vacío en la petición")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Tipo de integración requerido",
			Error:   "El parámetro 'type' no puede estar vacío",
		})
		return
	}

	var businessID *uint
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if id, err := strconv.ParseUint(businessIDStr, 10, 32); err == nil {
			id := uint(id)
			businessID = &id
		} else {
			h.logger.Error().Err(err).Str("business_id", businessIDStr).Str("type", integrationType).Msg("Error al parsear business_id en query parameter")
		}
	}

	integrationWithCreds, err := h.usecase.GetIntegrationByType(c.Request.Context(), integrationType, businessID)
	if err != nil {
		statusCode := http.StatusNotFound
		errorMsg := "Integración no encontrada"

		if errors.Is(err, domain.ErrIntegrationTypeNotFound) {
			errorMsg = "Tipo de integración no encontrado"
		} else if errors.Is(err, domain.ErrIntegrationCredentialsDecrypt) {
			statusCode = http.StatusInternalServerError
			errorMsg = "Error al procesar credenciales de la integración"
		} else if !errors.Is(err, domain.ErrIntegrationNotFound) {
			statusCode = http.StatusInternalServerError
			errorMsg = "Error interno del servidor al obtener la integración"
		}

		if businessID != nil {
			h.logger.Error().
				Err(err).
				Str("integration_type", integrationType).
				Uint("business_id", *businessID).
				Int("status_code", statusCode).
				Msg("Error al obtener integración por tipo y business_id en el usecase")
		} else {
			h.logger.Error().
				Err(err).
				Str("integration_type", integrationType).
				Int("status_code", statusCode).
				Msg("Error al obtener integración por tipo en el usecase")
		}
		c.JSON(statusCode, response.IntegrationErrorResponse{
			Success: false,
			Message: errorMsg,
			Error:   err.Error(),
		})
		return
	}

	// Retornar sin credenciales desencriptadas (por seguridad)
	integrationResp := mapper.ToIntegrationResponse(&integrationWithCreds.Integration)
	c.JSON(http.StatusOK, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración obtenida exitosamente",
		Data:    integrationResp,
	})
}
