package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/mappers"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/response"
)

// RegisterManualInvoice registra una factura externa asociada a una orden
func (h *handler) RegisterManualInvoice(c *gin.Context) {
	ctx := c.Request.Context()

	// Parsear request
	var req request.RegisterManualInvoice
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "invalid_request",
			Message: "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Resolver business_id (super admin o normal)
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, response.Error{
			Error:   "business_id_required",
			Message: "Se requiere seleccionar un negocio",
		})
		return
	}

	h.log.Info(ctx).
		Str("order_id", req.OrderID).
		Str("invoice_number", req.InvoiceNumber).
		Uint("business_id", businessID).
		Msg("Registrando factura manual")

	// Construir DTO
	dto := &dtos.RegisterManualInvoiceDTO{
		InvoiceNumber: req.InvoiceNumber,
		OrderID:       req.OrderID,
		BusinessID:    businessID,
	}

	// Llamar caso de uso
	invoice, err := h.useCase.RegisterManualInvoice(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Str("order_id", req.OrderID).Msg("Error registrando factura manual")
		handleDomainError(c, err, "manual_invoice_failed")
		return
	}

	// Convertir a response
	baseURL, bucket := h.getS3Config()
	resp := mappers.InvoiceToResponse(invoice, false, baseURL, bucket)

	h.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Msg("Factura manual registrada exitosamente")

	c.JSON(http.StatusCreated, resp)
}
