package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) RetryFailedInvoices(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "business_id es requerido",
		})
		return
	}

	queued, err := h.useCase.RetryFailedBulk(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"queued":  queued,
		"message": "Reconciliacion y reintento de fallidas en proceso",
	})
}
