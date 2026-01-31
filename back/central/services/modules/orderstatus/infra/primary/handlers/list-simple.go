package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/infra/primary/handlers/response"
)

// ListOrderStatusesSimple godoc
// @Summary      Listar estados de órdenes en formato simple
// @Description  Retorna solo ID, nombre, código e is_active de estados de órdenes para dropdowns/selectores
// @Tags         Order Statuses
// @Produce      json
// @Param        is_active  query     bool    false  "Filtrar por estado activo (default: true)"
// @Success      200        {object}  response.OrderStatusesSimpleResponse
// @Failure      500        {object}  map[string]interface{}
// @Router       /order-statuses/simple [get]
func (h *OrderStatusMappingHandlers) ListOrderStatusesSimple(c *gin.Context) {
	// Filtro is_active (por defecto true para solo activos)
	isActiveStr := c.DefaultQuery("is_active", "true")
	var isActive *bool
	if isActiveStr != "" {
		if isActiveValue, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &isActiveValue
		}
	}

	orderStatuses, err := h.uc.ListOrderStatuses(c.Request.Context(), isActive)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting order statuses for simple list")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener estados de órdenes",
			"error":   err.Error(),
		})
		return
	}

	// Mapear a formato simple
	simpleStatuses := make([]response.OrderStatusSimpleResponse, 0, len(orderStatuses))
	for _, status := range orderStatuses {
		// Si se filtra por is_active, todos los retornados tienen ese estado
		// Si no se filtra, asumimos que están activos (ya que es el comportamiento por defecto)
		statusIsActive := true
		if isActive != nil {
			statusIsActive = *isActive
		}

		simpleStatuses = append(simpleStatuses, response.OrderStatusSimpleResponse{
			ID:       status.ID,
			Name:     status.Name,
			Code:     status.Code,
			IsActive: statusIsActive,
		})
	}

	c.JSON(http.StatusOK, response.OrderStatusesSimpleResponse{
		Success: true,
		Message: "Estados de órdenes obtenidos exitosamente",
		Data:    simpleStatuses,
	})
}
