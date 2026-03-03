package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/internal/modules/orders/internal/infra/primary/handlers/mappers"
)

func (h *Handlers) GetReferenceData(c *gin.Context) {
	businessID := c.GetUint("testing_business_id")

	data, err := h.useCase.GetReferenceData(c.Request.Context(), businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": mappers.ReferenceDataToResponse(data)})
}
