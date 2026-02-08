package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/handlers/request"
)

// BulkCreateInvoices crea facturas masivamente de forma asíncrona usando RabbitMQ
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

	// Validar límite máximo
	if len(req.OrderIDs) > 500 {
		h.log.Warn(ctx).Int("order_count", len(req.OrderIDs)).Msg("Exceeded maximum bulk size")
		c.JSON(400, gin.H{"error": "Maximum 500 orders per batch"})
		return
	}

	h.log.Info(ctx).
		Int("order_count", len(req.OrderIDs)).
		Msg("Creating bulk invoice job")

	// Convertir a DTO de dominio
	dto := &dtos.BulkCreateInvoicesDTO{
		OrderIDs:   req.OrderIDs,
		BusinessID: req.BusinessID,
	}

	// Ejecutar caso de uso asíncrono - retorna jobID inmediatamente
	jobID, err := h.useCase.BulkCreateInvoicesAsync(ctx, dto)
	if err != nil {
		h.log.Error(ctx).Err(err).Msg("Failed to create bulk invoice job")
		c.JSON(500, gin.H{"error": "Failed to create bulk invoice job: " + err.Error()})
		return
	}

	h.log.Info(ctx).
		Str("job_id", jobID).
		Int("order_count", len(req.OrderIDs)).
		Msg("Bulk invoice job created successfully")

	// Retornar HTTP 202 Accepted con jobID
	c.JSON(202, gin.H{
		"job_id":       jobID,
		"status":       "processing",
		"total_orders": len(req.OrderIDs),
		"message":      "Invoices are being processed asynchronously",
		"status_url":   "/api/v1/invoicing/bulk-jobs/" + jobID,
	})
}
