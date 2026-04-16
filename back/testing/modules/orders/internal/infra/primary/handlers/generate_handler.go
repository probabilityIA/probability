package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/testing/modules/orders/internal/infra/primary/handlers/mappers"
)

func (h *Handlers) GenerateOrders(c *gin.Context) {
	businessID := c.GetUint("testing_business_id")

	var dto dtos.GenerateOrdersDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Forward the original token to the central API
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	result, err := h.useCase.GenerateOrders(c.Request.Context(), businessID, &dto, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    mappers.GenerateResultToResponse(result),
		"message": "Order generation completed",
	})
}
