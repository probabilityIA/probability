package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	invoicingErrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

func (uc *useCase) RegisterManualInvoice(ctx context.Context, dto *dtos.RegisterManualInvoiceDTO) (*entities.Invoice, error) {
	uc.log.Info(ctx).
		Str("order_id", dto.OrderID).
		Str("invoice_number", dto.InvoiceNumber).
		Uint("business_id", dto.BusinessID).
		Msg("Registrando factura manual externa")

	order, err := uc.repo.GetOrderByID(ctx, dto.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("order_id", dto.OrderID).Msg("Orden no encontrada")
		return nil, fmt.Errorf("orden no encontrada: %w", err)
	}

	exists, err := uc.repo.InvoiceExistsForOrder(ctx, order.ID, 0)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error verificando factura existente")
		return nil, fmt.Errorf("error verificando factura existente: %w", err)
	}
	if exists {
		return nil, invoicingErrors.ErrOrderAlreadyInvoiced
	}

	now := time.Now()
	invoice := &entities.Invoice{
		OrderID:       order.ID,
		BusinessID:    dto.BusinessID,
		InvoiceNumber: dto.InvoiceNumber,
		Status:        constants.InvoiceStatusIssued,
		IssuedAt:      &now,
		Subtotal:      order.Subtotal,
		Tax:           order.Tax,
		Discount:      order.Discount,
		ShippingCost:  dtos.InvoiceShippingCost(order),
		TotalAmount:   dtos.InvoiceTotalAmount(order),
		Currency:      order.Currency,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		CustomerPhone: order.CustomerPhone,
		CustomerDNI:   order.CustomerDNI,
	}

	notes := "Factura registrada manualmente (externa al sistema)"
	invoice.Notes = &notes

	if err := uc.repo.CreateInvoice(ctx, invoice); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error creando factura manual")
		return nil, fmt.Errorf("error creando factura manual: %w", err)
	}

	if err := uc.repo.UpdateOrderInvoiceInfo(ctx, order.ID, fmt.Sprintf("%d", invoice.ID), ""); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("No se pudo actualizar invoice_id en la orden (la factura fue creada)")
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Str("order_id", order.ID).
		Msg("Factura manual registrada exitosamente")

	if uc.ssePublisher != nil {
		invoice.OrderNumber = order.OrderNumber
		if err := uc.ssePublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Error publicando evento SSE de factura manual")
		}
	}

	return invoice, nil
}
