package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOrderStatuses godoc
// @Summary      Listar estados de órdenes de Probability
// @Description  Obtiene una lista de todos los estados de órdenes de Probability. Opcionalmente puede filtrar por estado activo/inactivo.
// @Tags         Order Statuses
// @Accept       json
// @Produce      json
// @Param        is_active  query     bool    false  "Filtrar por estado activo/inactivo (true=activos, false=inactivos, omitir=todos)"
// @Success      200        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]string
// @Router       /order-statuses [get]
func (h *OrderStatusMappingHandlers) ListOrderStatuses(c *gin.Context) {
	var isActive *bool

	// Filtro opcional por is_active
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActiveValue, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &isActiveValue
		}
	}

	result, err := h.uc.ListOrderStatuses(c.Request.Context(), isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estados de órdenes",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Estados de órdenes obtenidos exitosamente",
		"data":    result,
	})
}
