package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// GetInvoice obtiene una factura por ID
func (uc *useCase) GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error) {
<<<<<<< HEAD
	invoice, err := uc.invoiceRepo.GetByID(ctx, invoiceID)
=======
	invoice, err := uc.repo.GetInvoiceByID(ctx, invoiceID)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err != nil {
		return nil, errors.ErrInvoiceNotFound
	}
	return invoice, nil
}

<<<<<<< HEAD
// ListInvoices lista facturas con filtros
func (uc *useCase) ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error) {
	return uc.invoiceRepo.List(ctx, filters)
=======
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
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
}

// GetInvoicesByOrder obtiene todas las facturas de una orden
func (uc *useCase) GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error) {
<<<<<<< HEAD
	invoice, err := uc.invoiceRepo.GetByOrderID(ctx, orderID)
=======
	invoice, err := uc.repo.GetInvoiceByOrderID(ctx, orderID)
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err != nil {
		// Retornar lista vacía si no hay facturas
		return []*entities.Invoice{}, nil
	}
	return []*entities.Invoice{invoice}, nil
}
