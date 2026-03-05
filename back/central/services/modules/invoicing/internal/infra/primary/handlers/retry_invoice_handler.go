package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// RetryInvoice reintenta una factura fallida o consulta estado de una factura pending
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

	// Obtener factura para decidir qué acción tomar
	invoice, err := h.useCase.GetInvoice(ctx, uint(id))
	if err != nil {
		handleDomainError(c, err, "get_invoice_failed")
		return
	}

	// Decidir acción según estado:
	// - failed → RetryInvoice (re-envía POST con verificación de idempotencia)
	// - pending → CheckPendingInvoice (solo busca documento, NO re-envía POST)
	if invoice.Status == "pending" {
		h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Checking pending invoice status")
		err = h.useCase.CheckPendingInvoice(ctx, uint(id))
	} else {
		h.log.Info(ctx).Uint("invoice_id", uint(id)).Msg("Retrying failed invoice")
		err = h.useCase.RetryInvoice(ctx, uint(id))
	}

	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to process invoice")
		handleDomainError(c, err, "retry_invoice_failed")
		return
	}

	// Obtener factura actualizada
	invoice, err = h.useCase.GetInvoice(ctx, uint(id))
	if err != nil {
		handleDomainError(c, err, "get_invoice_failed")
		return
	}

	baseURL, bucket := h.getS3Config()
	resp := mappers.InvoiceToResponse(invoice, true, baseURL, bucket)

	c.JSON(http.StatusOK, resp)
}
