package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
)

// BulkCreateInvoices crea facturas masivamente a partir de una lista de order_ids
// POST /api/v1/invoicing/invoices/bulk
func (h *handler) BulkCreateInvoices(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var req request.BulkCreateInvoicesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error(ctx).Err(err).Msg("Invalid request body")
		c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	h.log.Info(ctx).
		Int("order_count", len(req.OrderIDs)).
		Msg("Bulk creating invoices")

	// Convertir a DTO de dominio
	dto := &dtos.BulkCreateInvoicesDTO{
		OrderIDs: req.OrderIDs,
	}

	// Ejecutar caso de uso
	result, err := h.useCase.BulkCreateInvoices(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to bulk create invoices")
		c.JSON(500, gin.H{"error": "Failed to create invoices: " + err.Error()})
		return
	}

	h.log.Info(ctx).
		Int("created", result.Created).
		Int("failed", result.Failed).
		Msg("Bulk invoice creation completed")

	// Retornar resultado (200 OK incluso si algunas fallaron)
	c.JSON(200, result)
}
