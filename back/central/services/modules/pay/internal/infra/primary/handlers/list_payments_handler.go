package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/infra/primary/handlers/response"
)

// ListPayments maneja GET /pay/transactions
func (h *handler) ListPayments(c *gin.Context) {
	ctx := c.Request.Context()

	businessID, ok := middleware.GetBusinessIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	txs, total, err := h.useCase.ListPayments(ctx, businessID, page, pageSize)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to list payments")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]*response.PaymentTransactionResponse, len(txs))
	for i, tx := range txs {
		data[i] = mappers.ToResponse(tx)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.PaginatedPaymentsResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
