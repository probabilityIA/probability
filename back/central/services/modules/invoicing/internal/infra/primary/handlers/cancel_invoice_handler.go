package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// CancelInvoice cancela una factura
func (h *handler) CancelInvoice(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener ID del path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid invoice ID",
		})
		return
	}

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Cancelling invoice")

	// Llamar caso de uso
	err = h.useCase.CancelInvoice(ctx, &dtos.CancelInvoiceDTO{InvoiceID: uint(id)})
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to cancel invoice")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "cancel_invoice_failed",
			Message: err.Error(),
		})
		return
	}

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Invoice cancelled successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Invoice cancelled successfully",
		"invoice_id": uint(id),
	})
}
