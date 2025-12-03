package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
)

// GetIntegrationByIDHandler obtiene una integración por su ID
//
//	@Summary		Obtener integración por ID
//	@Description	Obtiene una integración específica por su ID
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int		true	"ID de la integración"
//	@Success		200	{object}	response.IntegrationSuccessResponse
//	@Failure		400	{object}	response.IntegrationErrorResponse
//	@Failure		404	{object}	response.IntegrationErrorResponse
//	@Failure		500	{object}	response.IntegrationErrorResponse
//	@Router			/integrations/{id} [get]
func (h *IntegrationHandler) GetIntegrationByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "ID inválido",
			Error:   "El ID debe ser un número válido",
		})
		return
	}

	integration, err := h.usecase.GetIntegrationByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint64("id", id).Msg("Error al obtener integración")
		c.JSON(http.StatusNotFound, response.IntegrationErrorResponse{
			Success: false,
			Message: "Integración no encontrada",
			Error:   err.Error(),
		})
		return
	}

	integrationResp := mapper.ToIntegrationResponse(integration)
	c.JSON(http.StatusOK, response.IntegrationSuccessResponse{
		Success: true,
		Message: "Integración obtenida exitosamente",
		Data:    integrationResp,
	})
}
