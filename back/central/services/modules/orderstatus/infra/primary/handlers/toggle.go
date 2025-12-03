package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Toggle godoc
// @Summary      Activar/Desactivar mapeo
// @Description  Cambia el estado activo/inactivo de un mapeo
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      200  {object}  response.OrderStatusMappingResponse
// @Failure      400  {object}  map[string]string
// @Router       /order-status-mappings/{id}/toggle [patch]
func (h *OrderStatusMappingHandlers) Toggle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.uc.ToggleOrderStatusMappingActive(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toResponse(result))
}
