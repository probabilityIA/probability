package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// RetryInvoice reintenta una factura fallida
func (h *handler) RetryInvoice(c *gin.Context) {
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

	h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Retrying invoice")

	// Llamar caso de uso
	err = h.useCase.RetryInvoice(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to retry invoice")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "retry_invoice_failed",
			Message: err.Error(),
		})
		return
	}

	// Obtener factura actualizada
	invoice, err := h.useCase.GetInvoice(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to get retried invoice")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "get_invoice_failed",
			Message: err.Error(),
		})
		return
	}

	// Convertir a response
	resp := mappers.InvoiceToResponse(invoice, true)

	h.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("status", invoice.Status).
		Msg("Invoice retry completed")

	c.JSON(http.StatusOK, resp)
}
