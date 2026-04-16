package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// GenerateCashReceipt genera un recibo de caja para una factura ya emitida
func (h *handler) GenerateCashReceipt(c *gin.Context) {
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

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Generating cash receipt for invoice")

	if err := h.useCase.GenerateCashReceipt(ctx, uint(id)); err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to generate cash receipt")
		handleDomainError(c, err, "cash_receipt_failed")
		return
	}

	// Obtener factura actualizada
	invoice, err := h.useCase.GetInvoice(ctx, uint(id))
	if err != nil {
		handleDomainError(c, err, "get_invoice_failed")
		return
	}

	baseURL, bucket := h.getS3Config()
	resp := mappers.InvoiceToResponse(invoice, true, baseURL, bucket)

	c.JSON(http.StatusOK, resp)
}
