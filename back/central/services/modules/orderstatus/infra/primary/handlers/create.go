package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/request"
)

// Create godoc
// @Summary      Crear mapeo de estado de orden
// @Description  Crea un nuevo mapeo entre un estado de orden original y uno unificado
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateOrderStatusMappingRequest  true  "Datos del mapeo"
// @Success      201      {object}  response.OrderStatusMappingResponse
// @Failure      400      {object}  map[string]string
// @Router       /order-status-mappings [post]
func (h *OrderStatusMappingHandlers) Create(c *gin.Context) {
	var req request.CreateOrderStatusMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainEntity := toDomainCreate(&req)
	result, err := h.uc.CreateOrderStatusMapping(c.Request.Context(), domainEntity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toResponse(result))
}
