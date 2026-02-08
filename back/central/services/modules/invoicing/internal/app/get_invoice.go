package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// GetInvoice obtiene una factura por ID
func (uc *useCase) GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error) {
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, errors.ErrInvoiceNotFound
	}
	return invoice, nil
}

// ListInvoices lista facturas con filtros y paginación
func (uc *useCase) ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error) {
	return uc.repo.ListInvoices(ctx, filters)
}

// GetInvoiceSyncLogs obtiene los logs de sincronización de una factura
func (uc *useCase) GetInvoiceSyncLogs(ctx context.Context, invoiceID uint) ([]*entities.InvoiceSyncLog, error) {
	// Verificar que la factura existe
	_, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, errors.ErrInvoiceNotFound
	}

	return uc.repo.GetSyncLogsByInvoiceID(ctx, invoiceID)
}

// GetInvoicesByOrder obtiene todas las facturas de una orden
func (uc *useCase) GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error) {
	invoice, err := uc.repo.GetInvoiceByOrderID(ctx, orderID)
	if err != nil {
		// Retornar lista vacía si no hay facturas
		return []*entities.Invoice{}, nil
	}
	return []*entities.Invoice{invoice}, nil
}
