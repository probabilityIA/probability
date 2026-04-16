package handlerintegrations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/request"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
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
//	@Param			integration_type_id		query		int		false	"Filtrar por ID del tipo de integración"
//	@Param			integration_type_code	query		string	false	"Filtrar por código del tipo de integración"
//	@Param			category	query		string	false	"Filtrar por categoría"
//	@Param			business_id	query		int		false	"Filtrar por business ID"
//	@Param			is_active	query		bool	false	"Filtrar por estado activo"
//	@Param			search		query		string	false	"Buscar por nombre o código"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]interface{}
//	@Failure		401			{object}	map[string]interface{}
//	@Failure		500			{object}	map[string]interface{}
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

	// Seguridad: Filtrar por business_id si no es super admin o platform
	if businessID, exists := c.Get("business_id"); exists {
		if bID, ok := businessID.(uint); ok && bID > 0 {
			// Si es un usuario de negocio, forzar el filtrado por su business_id
			req.BusinessID = &bID
		}
	} else {
		// Fallback: intentar ver si es super admin por el flag explícito
		if isSuperAdmin, exists := c.Get("is_super_admin"); exists {
			if isSuper, ok := isSuperAdmin.(bool); ok && !isSuper {
				// Si NO es super admin y no se pudo obtener business_id (caso raro),
				// podríamos loguear error o negar acceso, pero por ahora confiamos en el middleware.
				// Ojo: si business_id no está seteado en contexto, asume 0/nil en req solo si así vino.
			}
		}
	}

	filters := mapper.ToIntegrationFilters(req)
	integrations, total, err := h.usecase.ListIntegrations(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error().
			Err(err).
			Int("page", filters.Page).
			Int("page_size", filters.PageSize).
			Int("status_code", http.StatusInternalServerError).
			Msg("Error al listar integraciones en el usecase")
		c.JSON(http.StatusInternalServerError, response.IntegrationErrorResponse{
			Success: false,
			Message: "Error interno del servidor al obtener la lista de integraciones",
			Error:   err.Error(),
		})
		return
	}

	// Obtener URL base de S3 para construir URLs completas
	imageURLBase := h.getImageURLBase()
	response := mapper.ToIntegrationListResponse(integrations, total, filters.Page, filters.PageSize, imageURLBase)
	c.JSON(http.StatusOK, response)
}
