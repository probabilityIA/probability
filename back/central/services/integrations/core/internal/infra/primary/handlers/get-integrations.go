package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/response"
)

// GetIntegrationsHandler obtiene la lista de integraciones
//
//	@Summary		Obtener integraciones
//	@Description	Obtiene una lista paginada de integraciones con filtros
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int		false	"Número de página"	default(1)
//	@Param			page_size	query		int		false	"Tamaño de página"	default(10)
//	@Param			type		query		string	false	"Filtrar por tipo"
//	@Param			category	query		string	false	"Filtrar por categoría"
//	@Param			business_id	query		int		false	"Filtrar por business ID"
//	@Param			is_active	query		bool	false	"Filtrar por estado activo"
//	@Param			search		query		string	false	"Buscar por nombre o código"
//	@Success		200			{object}	response.IntegrationListResponse
//	@Failure		400			{object}	response.IntegrationErrorResponse
//	@Failure		401			{object}	response.IntegrationErrorResponse
//	@Failure		500			{object}	response.IntegrationErrorResponse
//	@Router			/integrations [get]
func (h *IntegrationHandler) GetIntegrationsHandler(c *gin.Context) {
	var req request.GetIntegrationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al parsear parámetros de consulta")
		c.JSON(http.StatusBadRequest, response.IntegrationErrorResponse{
			Success: false,
			Message: "Parámetros de consulta inválidos",
			Error:   err.Error(),
		})
		return
	}

	filters := mapper.ToIntegrationFilters(req)
	integrations, total, err := h.usecase.ListIntegrations(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al listar integraciones")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error interno del servidor",
			Error:   err.Error(),
		})
		return
	}

	response := mapper.ToIntegrationListResponse(integrations, total, filters.Page, filters.PageSize)
	c.JSON(http.StatusOK, response)
}
