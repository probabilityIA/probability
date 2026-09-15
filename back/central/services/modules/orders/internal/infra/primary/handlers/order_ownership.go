package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	orderErrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
)

func (h *Handlers) ensureOrderOwnership(c *gin.Context, orderID string) bool {
	tokenBusinessID, _ := c.Get("business_id")
	requesterBusinessID, _ := tokenBusinessID.(uint)
	if requesterBusinessID == 0 {
		return true
	}

	order, err := h.orderCRUD.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, orderErrors.ErrOrderNotFound) || strings.Contains(err.Error(), "order not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Orden no encontrada",
				"error":   "order not found",
			})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al validar la orden",
			"error":   err.Error(),
		})
		return false
	}

	if order == nil || order.BusinessID == nil || *order.BusinessID != requesterBusinessID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "No autorizado",
			"error":   "la orden no pertenece a este negocio",
		})
		return false
	}

	return true
}
