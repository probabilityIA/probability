package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

type markInvoiceAsCancelledRequest struct {
	Reason string `json:"reason"`
}

func (h *handler) MarkInvoiceAsCancelled(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid invoice ID",
		})
		return
	}

	var req markInvoiceAsCancelledRequest
	_ = c.ShouldBindJSON(&req)
	if req.Reason == "" {
		req.Reason = "Marked as cancelled via admin endpoint"
	}

	userID := c.GetUint("user_id")

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Str("reason", req.Reason).Msg("Marking invoice as cancelled")

	err = h.useCase.MarkInvoiceAsCancelled(ctx, &dtos.MarkInvoiceAsCancelledDTO{
		InvoiceID:         uint(id),
		Reason:            req.Reason,
		CancelledByUserID: userID,
	})
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to mark invoice as cancelled")
		handleDomainError(c, err, "mark_invoice_as_cancelled_failed")
		return
	}

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Invoice marked as cancelled successfully")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Factura marcada como cancelada",
		"invoice_id": uint(id),
	})
}
