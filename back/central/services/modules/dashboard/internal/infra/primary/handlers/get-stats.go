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
	// Obtener business_id del contexto
	var businessID *uint
	if businessIDCtx, exists := c.Get("business_id"); exists {
		if bID, ok := businessIDCtx.(uint); ok {
			// Si business_id == 0, es super user (ver todo)
			// Si business_id > 0, filtrar por ese negocio
			if bID > 0 {
				businessID = &bID
			}
			// Si bID == 0, businessID queda nil (super user ve todo)
		}
	}

	// Si es super admin (businessID == nil), permitir filtrar por business_id opcional del query parameter
	// Esto debe estar FUERA del bloque anterior para que funcione cuando no hay business_id en contexto
	if businessID == nil {
		if businessIDParam := c.Query("business_id"); businessIDParam != "" {
			if parsedID, err := strconv.ParseUint(businessIDParam, 10, 32); err == nil && parsedID > 0 {
				filteredID := uint(parsedID)
				businessID = &filteredID
			}
		}
	}

	// Obtener estadísticas del caso de uso
	stats, err := h.uc.GetDashboardStats(c.Request.Context(), businessID)
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
