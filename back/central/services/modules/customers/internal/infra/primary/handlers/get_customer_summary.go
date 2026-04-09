package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

func (h *Handlers) GetCustomerSummary(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || customerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	summary, err := h.uc.GetCustomerSummary(c.Request.Context(), businessID, uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if summary == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer summary not found"})
		return
	}

	c.JSON(http.StatusOK, response.SummaryFromEntity(summary))
}
