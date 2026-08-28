package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

func (h *handler) RefreshInvoice(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_id",
			Message: "Invalid invoice ID",
		})
		return
	}

	if err := h.useCase.RefreshInvoiceFromProvider(ctx, uint(id)); err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to refresh invoice from provider")
		handleDomainError(c, err, "refresh_invoice_failed")
		return
	}

	invoice, err := h.useCase.GetInvoice(ctx, uint(id))
	if err != nil {
		handleDomainError(c, err, "get_invoice_failed")
		return
	}

	baseURL, bucket := h.getS3Config()
	c.JSON(http.StatusAccepted, mappers.InvoiceToResponse(invoice, true, baseURL, bucket))
}
