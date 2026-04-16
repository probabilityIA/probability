package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTopSellingDays obtiene los TOP N días de mayor demanda
// @Summary      Obtener TOP días de mayor demanda
// @Description  Retorna los días con más órdenes creadas en toda la historia del negocio
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Param        business_id  query    int     false  "ID del business para filtrar (solo super admin)"
// @Param        limit        query    int     false  "Número máximo de días a retornar (default: 5)"
// @Success      200  {array}   domain.TopSellingDay
// @Failure      500  {object}  map[string]interface{}
// @Router       /dashboard/top-selling-days [get]
func (h *DashboardHandlers) GetTopSellingDays(c *gin.Context) {
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
		if bID, ok := businessIDCtx.(uint); ok && bID > 0 {
			businessID = &bID
		} else if bIDFloat, ok := businessIDCtx.(float64); ok && bIDFloat > 0 {
			uID := uint(bIDFloat)
			businessID = &uID
		}
	}

	// 2. Si es super admin, permitir override por query param o ver todo (nil)
	if isSuperAdmin {
		if businessIDParam := c.Query("business_id"); businessIDParam != "" {
			if parsedID, err := strconv.ParseUint(businessIDParam, 10, 32); err == nil && parsedID > 0 {
				filteredID := uint(parsedID)
				businessID = &filteredID
			}
		}
	} else {
		// 3. Si NO es super admin, businessID DEBE estar seteado
		if businessID == nil {
			h.logger.Error().Msg("Intento de acceso a top días sin business_id en contexto para usuario no admin")
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado"})
			return
		}
	}

	// Obtener parámetro limit (default: 5)
	limit := 5
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Obtener integration_id si está especificado
	var integrationID *uint
	if integrationIDParam := c.Query("integration_id"); integrationIDParam != "" {
		if parsedID, err := strconv.ParseUint(integrationIDParam, 10, 32); err == nil && parsedID > 0 {
			filteredID := uint(parsedID)
			integrationID = &filteredID
		}
	}

	// Obtener datos del caso de uso
	topDays, err := h.uc.GetTopSellingDays(c.Request.Context(), businessID, integrationID, limit)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener TOP días de mayor demanda")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener TOP días de mayor demanda",
			"error":   err.Error(),
		})
		return
	}

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "TOP días obtenidos exitosamente",
		"data":    topDays,
	})
}
