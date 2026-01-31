package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// CreateInvoice crea una factura manualmente
func (h *handler) CreateInvoice(c *gin.Context) {
	ctx := c.Request.Context()

	// Parsear request
	var req request.CreateInvoice
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	h.log.Info(ctx).
		Str("order_id", req.OrderID).
		Msg("Creating invoice manually")

	// Convertir a DTO de dominio
	dto := mappers.CreateInvoiceRequestToDTO(&req)

	// Llamar caso de uso
	invoice, err := h.useCase.CreateInvoice(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Str("order_id", req.OrderID).Msg("Failed to create invoice")
		c.JSON(http.StatusInternalServerError, response.Error{
			Error:   "invoice_creation_failed",
			Message: err.Error(),
		})
		return
	}

	// Convertir a response
	resp := mappers.InvoiceToResponse(invoice, true) // Incluir items

	h.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Msg("Invoice created successfully")

	c.JSON(http.StatusCreated, resp)
}
