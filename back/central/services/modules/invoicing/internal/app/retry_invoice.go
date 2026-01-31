package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// RetryInvoice reintenta la creación de una factura fallida
func (uc *useCase) RetryInvoice(ctx context.Context, invoiceID uint) error {
	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Retrying invoice creation")

	// 1. Obtener factura
	invoice, err := uc.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return errors.ErrInvoiceNotFound
	}

	// 2. Validar que esté en estado failed
	if invoice.Status != constants.InvoiceStatusFailed {
		return errors.ErrRetryNotAllowed
	}

	// 3. Obtener logs de sincronización
	logs, err := uc.syncLogRepo.GetByInvoiceID(ctx, invoiceID)
	if err != nil || len(logs) == 0 {
		return fmt.Errorf("no sync logs found for invoice")
	}

	lastLog := logs[len(logs)-1]

	// 4. Validar que no se haya excedido el máximo de reintentos
	if lastLog.RetryCount >= lastLog.MaxRetries {
		return errors.ErrMaxRetriesExceeded
	}

	// 5. Reintentar creación usando CreateInvoice
	dto := &dtos.CreateInvoiceDTO{
		OrderID:  invoice.OrderID,
		Notes:    invoice.Notes,
		IsManual: true, // Los reintentos se consideran manuales
	}

	_, err = uc.CreateInvoice(ctx, dto)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Retry failed")
		return err
	}

	uc.log.Info(ctx).Uint("invoice_id", invoiceID).Msg("Invoice retry successful")
	return nil
}
