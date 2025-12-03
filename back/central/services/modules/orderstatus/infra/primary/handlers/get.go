package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get godoc
// @Summary      Obtener mapeo de estado de orden
// @Description  Obtiene un mapeo de estado de orden por su ID
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del mapeo"
// @Success      200  {object}  response.OrderStatusMappingResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /order-status-mappings/{id} [get]
func (h *OrderStatusMappingHandlers) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.uc.GetOrderStatusMapping(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toResponse(result))
}
