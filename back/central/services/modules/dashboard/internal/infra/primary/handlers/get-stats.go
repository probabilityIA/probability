package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
)

// GetStats obtiene las estadísticas del dashboard
// @Summary      Obtener estadísticas del dashboard
// @Description  Retorna estadísticas agregadas de órdenes, transportadores, productos, envíos y businesses (si es super admin)
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Param        business_id  query    int     false  "ID del business para filtrar (solo super admin)"
// @Success      200  {object}  domain.DashboardStatsResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /dashboard/stats [get]
func (h *DashboardHandlers) GetStats(c *gin.Context) {
	// Inicializar variables
	var businessID *uint
	isSuperAdmin := false

	// Verificar si es super admin
	if isSuperAdminCtx, exists := c.Get("is_super_admin"); exists {
		if isSuper, ok := isSuperAdminCtx.(bool); ok {
			isSuperAdmin = isSuper
		}
	}

	// 1. Intentar obtener business_id del contexto (para usuarios normales)
	if businessIDCtx, exists := c.Get("business_id"); exists {
		// Log para debug
		h.logger.Info().Interface("raw_context_val", businessIDCtx).Msg("Debug: Checking business_id in context")

		if bID, ok := businessIDCtx.(uint); ok && bID > 0 {
			businessID = &bID
		} else if bIDFloat, ok := businessIDCtx.(float64); ok && bIDFloat > 0 {
			// Handle float64 case just in case
			uID := uint(bIDFloat)
			businessID = &uID
		}
	}

	h.logger.Info().
		Bool("is_super_admin", isSuperAdmin).
		Interface("final_business_id", businessID).
		Msg("Debug: GetStats Filtering Decision")

	// 2. Si es super admin, permitir override por query param o ver todo (nil)
	if isSuperAdmin {
		if businessIDParam := c.Query("business_id"); businessIDParam != "" {
			if parsedID, err := strconv.ParseUint(businessIDParam, 10, 32); err == nil && parsedID > 0 {
				filteredID := uint(parsedID)
				businessID = &filteredID
			}
		}
		// Si no hay query param, businessID sigue siendo nil o lo que tenga el contexto
	} else {
		// 3. Si NO es super admin, businessID DEBE estar seteado.
		if businessID == nil {
			h.logger.Error().Msg("Intento de acceso a estadísticas sin business_id en contexto para usuario no admin")
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado"})
			return
		}
	}

	var integrationID *uint
	if integrationIDParam := c.Query("integration_id"); integrationIDParam != "" {
		if parsedID, err := strconv.ParseUint(integrationIDParam, 10, 32); err == nil && parsedID > 0 {
			filteredID := uint(parsedID)
			integrationID = &filteredID
		}
	}

	// Obtener estadísticas del caso de uso
	stats, err := h.uc.GetDashboardStats(c.Request.Context(), businessID, integrationID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener estadísticas del dashboard")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estadísticas del dashboard",
			"error":   err.Error(),
		})
		return
	}

	// Retornar respuesta exitosa
	response := domain.DashboardStatsResponse{
		Success: true,
		Message: "Estadísticas obtenidas exitosamente",
		Data:    *stats,
	}

	c.JSON(http.StatusOK, response)
}
