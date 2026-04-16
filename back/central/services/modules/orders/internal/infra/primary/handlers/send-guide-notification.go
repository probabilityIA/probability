package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendGuideNotification envia la notificacion de guia de envio por WhatsApp
// @Summary Enviar notificacion de guia por WhatsApp
// @Description Publica un evento para enviar la notificacion de guia de envio via WhatsApp
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Guide notification sent"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /orders/{id}/send-guide-notification [post]
func (h *Handlers) SendGuideNotification(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order_id is required",
		})
		return
	}

	if err := h.sendGuideNotificationUC.SendGuideNotification(c.Request.Context(), orderID); err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if errorMsg == "error getting order: order not found" {
			statusCode = http.StatusNotFound
		} else if errorMsg == "order does not have customer phone" ||
			errorMsg == "order does not have tracking number" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "guide_notification_sent",
		"order_id": orderID,
	})
}
