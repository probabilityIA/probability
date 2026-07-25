package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/mappers"
)

func (h *Handlers) GetOrderNotifications(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "order_id is required"})
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "business_id es requerido para super admin"})
		return
	}

	result, err := h.orderCRUD.GetOrderNotifications(c.Request.Context(), orderID, businessID)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "order not found":
			status = http.StatusNotFound
		case "la orden no pertenece al negocio":
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.OrderNotificationsToResponse(result))
}
