package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/mappers"
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
func (h *handler) ListOrderStatusesSimple(c *gin.Context) {
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

	// Convertir a formato simple usando mapper
	c.JSON(http.StatusOK, mappers.StatusInfoListToSimpleResponse(orderStatuses))
}
