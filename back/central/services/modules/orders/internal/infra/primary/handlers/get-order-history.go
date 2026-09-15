package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetOrderHistory(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de orden inválido",
			"error":   "El ID de la orden es requerido",
		})
		return
	}

	if !h.ensureOrderOwnership(c, id) {
		return
	}

	history, err := h.orderCRUD.GetOrderHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener historial",
			"error":   err.Error(),
		})
		return
	}

	httpHistory := make([]response.OrderHistoryResponse, len(history))
	for i, h := range history {
		httpHistory[i] = response.OrderHistoryResponse{
			ID:             h.ID,
			CreatedAt:      h.CreatedAt,
			OrderID:        h.OrderID,
			PreviousStatus: h.PreviousStatus,
			NewStatus:      h.NewStatus,
			ChangedBy:      h.ChangedBy,
			Source:         h.Source,
			ChangedByName:  h.ChangedByName,
			Reason:         h.Reason,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Historial obtenido exitosamente",
		"data":    httpHistory,
	})
}
