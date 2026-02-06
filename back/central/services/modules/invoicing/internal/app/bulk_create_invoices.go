package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
)

// BulkCreateInvoices crea múltiples facturas iterando sobre las órdenes proporcionadas
// Cada factura se procesa individualmente usando CreateInvoice, por lo que:
// - Se aplican todas las validaciones existentes
// - Se integra con Softpymes/proveedor configurado
// - Se publican eventos individuales
// - Si una factura falla, las demás continúan procesándose
func (uc *useCase) BulkCreateInvoices(ctx context.Context, dto *dtos.BulkCreateInvoicesDTO) (*dtos.BulkCreateResult, error) {
	result := &dtos.BulkCreateResult{
		Created: 0,
		Failed:  0,
		Results: make([]dtos.BulkInvoiceResult, 0, len(dto.OrderIDs)),
	}

	uc.log.Info(ctx).
		Int("total_orders", len(dto.OrderIDs)).
		Msg("Starting bulk invoice creation")

	// Procesar cada orden individualmente
	for _, orderID := range dto.OrderIDs {
		createDTO := &dtos.CreateInvoiceDTO{
			OrderID:  orderID,
			IsManual: true, // Marca como manual para logging/auditoría
		}

		invoice, err := uc.CreateInvoice(ctx, createDTO)

		if err != nil {
			result.Failed++
			errMsg := err.Error()
			result.Results = append(result.Results, dtos.BulkInvoiceResult{
				OrderID: orderID,
				Success: false,
				Error:   &errMsg,
			})
			uc.log.Warn(ctx).
				Err(err).
				Str("order_id", orderID).
				Msg("Failed to create invoice in bulk")
		} else {
			result.Created++
			result.Results = append(result.Results, dtos.BulkInvoiceResult{
				OrderID:   orderID,
				Success:   true,
				InvoiceID: &invoice.ID,
			})
			uc.log.Info(ctx).
				Str("order_id", orderID).
				Uint("invoice_id", invoice.ID).
				Msg("Invoice created successfully in bulk")
		}
	}

	uc.log.Info(ctx).
		Int("created", result.Created).
		Int("failed", result.Failed).
		Msg("Bulk invoice creation completed")

	return result, nil
}
