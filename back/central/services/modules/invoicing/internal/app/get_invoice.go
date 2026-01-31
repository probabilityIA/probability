package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

// GetInvoice obtiene una factura por ID
func (uc *useCase) GetInvoice(ctx context.Context, invoiceID uint) (*entities.Invoice, error) {
	invoice, err := uc.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, errors.ErrInvoiceNotFound
	}
	return invoice, nil
}

// ListInvoices lista facturas con filtros
func (uc *useCase) ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error) {
	return uc.invoiceRepo.List(ctx, filters)
}

// GetInvoicesByOrder obtiene todas las facturas de una orden
func (uc *useCase) GetInvoicesByOrder(ctx context.Context, orderID string) ([]*entities.Invoice, error) {
	invoice, err := uc.invoiceRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		// Retornar lista vacía si no hay facturas
		return []*entities.Invoice{}, nil
	}
	return []*entities.Invoice{invoice}, nil
}
