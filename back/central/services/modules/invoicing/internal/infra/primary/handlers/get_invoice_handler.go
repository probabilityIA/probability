package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// GetInvoice obtiene una factura por ID
func (h *handler) GetInvoice(c *gin.Context) {
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

	h.log.Debug(ctx).Uint("invoice_id", uint(id)).Msg("Getting invoice")

	// Llamar caso de uso
	invoice, err := h.useCase.GetInvoice(ctx, uint(id))
	if err != nil {
		h.log.Error(ctx).Err(err).Uint("invoice_id", uint(id)).Msg("Failed to get invoice")
		c.JSON(http.StatusNotFound, response.Error{
			Error:   "invoice_not_found",
			Message: err.Error(),
		})
		return
	}

	// Convertir a response
	resp := mappers.InvoiceToResponse(invoice, true) // Incluir items

	c.JSON(http.StatusOK, resp)
}
