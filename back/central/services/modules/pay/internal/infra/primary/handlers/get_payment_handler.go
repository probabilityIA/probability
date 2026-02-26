package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
)

// GetPayment maneja GET /pay/transactions/:id
func (h *handler) GetPayment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tx, err := h.useCase.GetPayment(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint64("id", id).Msg("Failed to get payment")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mappers.ToResponse(tx))
}
