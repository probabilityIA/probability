package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
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
//	@Success		200			{object}	response.IntegrationSuccessResponse
//	@Failure		400			{object}	response.IntegrationErrorResponse
//	@Failure		404			{object}	response.IntegrationErrorResponse
//	@Failure		500			{object}	response.IntegrationErrorResponse
//	@Router			/integrations/type/{type} [get]
func (h *IntegrationHandler) GetIntegrationByTypeHandler(c *gin.Context) {
	integrationType := c.Param("type")
	if integrationType == "" {
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Tipo de integración requerido",
		})
		return
	}

	var businessID *uint
	if businessIDStr := c.Query("business_id"); businessIDStr != "" {
		if id, err := strconv.ParseUint(businessIDStr, 10, 32); err == nil {
			id := uint(id)
			businessID = &id
		}
	}

	integrationWithCreds, err := h.usecase.GetIntegrationByType(c.Request.Context(), integrationType, businessID)
	if err != nil {
		h.logger.Error().Err(err).Str("type", integrationType).Msg("Error al obtener integración por tipo")
		c.JSON(http.StatusNotFound, response.IntegrationErrorResponse{
			Success: false,
			Message: "Integración no encontrada",
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
