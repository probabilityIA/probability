package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestConfirmation solicita la confirmación de una orden
// @Summary Solicitar confirmación de orden
// @Description Publica un evento para solicitar confirmación de una orden vía WhatsApp
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Confirmation request sent"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/{id}/request-confirmation [post]
func (h *Handlers) RequestConfirmation(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order_id is required",
		})
		return
	}

	// Llamar al caso de uso
	if err := h.requestConfirmationUC.RequestConfirmation(c.Request.Context(), orderID); err != nil {
		// Determinar código de error apropiado
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		// Errores específicos
		if errorMsg == "error getting order: order not found" {
			statusCode = http.StatusNotFound
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
