package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/request"
)

// Update godoc
// @Summary      Actualizar mapeo de estado de orden
// @Description  Actualiza un mapeo existente
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        id       path      int                                     true  "ID del mapeo"
// @Param        request  body      request.UpdateOrderStatusMappingRequest  true  "Datos a actualizar"
// @Success      200      {object}  response.OrderStatusMappingResponse
// @Failure      400      {object}  map[string]string
// @Router       /order-status-mappings/{id} [put]
func (h *OrderStatusMappingHandlers) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.UpdateOrderStatusMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainEntity := toDomainUpdate(&req)
	result, err := h.uc.UpdateOrderStatusMapping(c.Request.Context(), uint(id), domainEntity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toResponse(result))
}
