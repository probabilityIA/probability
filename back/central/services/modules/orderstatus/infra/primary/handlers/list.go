package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// List godoc
// @Summary      Listar mapeos de estado de orden
// @Description  Obtiene una lista de mapeos de estado de orden con filtros opcionales
// @Tags         Order Status Mappings
// @Accept       json
// @Produce      json
// @Param        integration_type  query     string  false  "Filtrar por tipo de integración (shopify, whatsapp, mercadolibre)"
// @Success      200               {object}  response.OrderStatusMappingsListResponse
// @Failure      500               {object}  map[string]string
// @Router       /order-status-mappings [get]
func (h *OrderStatusMappingHandlers) List(c *gin.Context) {
	filters := make(map[string]interface{})
	if integrationType := c.Query("integration_type"); integrationType != "" {
		filters["integration_type"] = integrationType
	}

	result, total, err := h.uc.ListOrderStatusMappings(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toListResponse(result, total))
}
