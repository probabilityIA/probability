package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) RequestConfirmation(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order_id is required",
		})
		return
	}

	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id es requerido para super admin"})
		return
	}

	if err := h.requestConfirmationUC.RequestConfirmation(c.Request.Context(), orderID, businessID); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if errorMsg == "error getting order: order not found" {
			statusCode = http.StatusNotFound
		} else if errorMsg == "la orden no pertenece al negocio" {
			statusCode = http.StatusForbidden
		} else if errorMsg == "order does not have customer phone" ||
			errorMsg == "order is already confirmed" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "confirmation_requested",
		"order_id": orderID,
	})
}
