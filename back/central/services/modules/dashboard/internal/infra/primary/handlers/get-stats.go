package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/dashboard/internal/domain"
)

// GetStats obtiene las estadísticas del dashboard
// @Summary      Obtener estadísticas del dashboard
// @Description  Retorna estadísticas agregadas de órdenes (total, por tipo de integración, top clientes, por ubicación)
// @Tags         Dashboard
// @Accept       json
// @Produce      json
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
