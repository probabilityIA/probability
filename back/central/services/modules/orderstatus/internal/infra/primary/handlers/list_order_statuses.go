package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/infra/primary/handlers/response"
)

// ListOrderStatuses godoc
// @Summary      Listar estados de órdenes de Probability
// @Description  Obtiene una lista de todos los estados de órdenes de Probability. Opcionalmente puede filtrar por estado activo/inactivo.
// @Tags         Order Statuses
// @Accept       json
// @Produce      json
// @Param        is_active  query     bool    false  "Filtrar por estado activo/inactivo (true=activos, false=inactivos, omitir=todos)"
// @Success      200        {object}  response.OrderStatusesCatalogResponse
// @Failure      500        {object}  map[string]string
// @Router       /order-statuses [get]
func (h *handler) ListOrderStatuses(c *gin.Context) {
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

	// Mapear a response con JSON tags lowercase
	data := make([]response.OrderStatusCatalogResponse, len(result))
	for i, s := range result {
		data[i] = response.OrderStatusCatalogResponse{
			ID:          s.ID,
			Code:        s.Code,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			Color:       s.Color,
			Priority:    s.Priority,
			IsActive:    s.IsActive,
		}
	}

	c.JSON(http.StatusOK, response.OrderStatusesCatalogResponse{
		Success: true,
		Message: "Estados de órdenes obtenidos exitosamente",
		Data:    data,
	})
}
