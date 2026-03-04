package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) DeleteAllOrders(c *gin.Context) {
	businessID := c.GetUint("testing_business_id")

	count, err := h.useCase.DeleteAllOrders(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": count,
		"message": fmt.Sprintf("%d órdenes eliminadas correctamente", count),
	})
}
