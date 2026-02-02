package handlerintegrations

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/response"
)

// GetIntegrationsSimpleHandler obtiene lista simple de integraciones
//
//	@Summary		Obtener lista simple de integraciones
//	@Description	Retorna solo ID, nombre, tipo, business_id e is_active de integraciones para dropdowns/selectores
//	@Tags			Integrations
//	@Produce		json
//	@Security		BearerAuth
//	@Param			business_id	query		int		false	"Filtrar por business ID"
//	@Param			is_active	query		bool	false	"Filtrar por estado activo (default: true)"
//	@Success		200			{object}	response.GetIntegrationsSimpleResponse
//	@Failure		500			{object}	map[string]interface{}
//	@Router			/integrations/simple [get]
func (h *IntegrationHandler) GetIntegrationsSimpleHandler(c *gin.Context) {
	// Parsear filtros opcionales
	var businessID *uint
	if bid := c.Query("business_id"); bid != "" {
		if id, err := strconv.ParseUint(bid, 10, 32); err == nil {
			uidValue := uint(id)
			businessID = &uidValue
		}
	}

	// Filtro is_active (por defecto true para solo activas)
	isActiveStr := c.DefaultQuery("is_active", "true")
	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	// Seguridad: Filtrar por business_id del usuario si no es super admin
	if userBusinessID, exists := c.Get("business_id"); exists {
		if bID, ok := userBusinessID.(uint); ok && bID > 0 {
			businessID = &bID
		}
	}

	// Crear filtros para el use case
	filters := domain.IntegrationFilters{
		Page:       1,
		PageSize:   1000, // Obtener todas las integraciones (suficiente para selectores)
		BusinessID: businessID,
		IsActive:   isActive,
	}

	integrations, _, err := h.usecase.ListIntegrations(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error getting integrations for simple list")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener integraciones",
			"error":   err.Error(),
		})
		return
	}

	// Mapear a formato simple
	simpleIntegrations := make([]response.IntegrationSimpleResponse, 0, len(integrations))
	for _, integration := range integrations {
		// Obtener el código del tipo de integración
		typeCode := ""
		categoryCode := ""
		categoryName := ""
		categoryColor := ""

		if integration.IntegrationType != nil {
			typeCode = integration.IntegrationType.Code

			// Extraer categoría desde IntegrationType.Category
			if integration.IntegrationType.Category != nil {
				categoryCode = integration.IntegrationType.Category.Code
				categoryName = integration.IntegrationType.Category.Name
				categoryColor = integration.IntegrationType.Category.Color
			}
		}

		simpleIntegrations = append(simpleIntegrations, response.IntegrationSimpleResponse{
			ID:            integration.ID,
			Name:          integration.Name,
			Type:          typeCode,
			Category:      categoryCode,
			CategoryName:  categoryName,
			CategoryColor: categoryColor,
			BusinessID:    integration.BusinessID,
			IsActive:      integration.IsActive,
		})
	}

	c.JSON(http.StatusOK, response.GetIntegrationsSimpleResponse{
		Success: true,
		Message: "Integraciones obtenidas exitosamente",
		Data:    simpleIntegrations,
	})
}
