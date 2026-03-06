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

// RegisterManualInvoice registra una factura externa (hecha por fuera del sistema).
// No envía nada al proveedor de facturación: solo crea el registro en la BD
// con status="issued" y actualiza la orden para que no se vuelva a facturar.
func (uc *useCase) RegisterManualInvoice(ctx context.Context, dto *dtos.RegisterManualInvoiceDTO) (*entities.Invoice, error) {
	uc.log.Info(ctx).
		Str("order_id", dto.OrderID).
		Str("invoice_number", dto.InvoiceNumber).
		Uint("business_id", dto.BusinessID).
		Msg("Registrando factura manual externa")

	// 1. Obtener la orden
	order, err := uc.repo.GetOrderByID(ctx, dto.OrderID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("order_id", dto.OrderID).Msg("Orden no encontrada")
		return nil, fmt.Errorf("orden no encontrada: %w", err)
	}

	// 2. Verificar que no exista ya una factura para esta orden
	exists, err := uc.repo.InvoiceExistsForOrder(ctx, order.ID, 0)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error verificando factura existente")
		return nil, fmt.Errorf("error verificando factura existente: %w", err)
	}
	if exists {
		return nil, invoicingErrors.ErrOrderAlreadyInvoiced
	}

	// 3. Crear la factura con status "issued" (ya fue emitida externamente)
	now := time.Now()
	invoice := &entities.Invoice{
		OrderID:       order.ID,
		BusinessID:    dto.BusinessID,
		InvoiceNumber: dto.InvoiceNumber,
		Status:        constants.InvoiceStatusIssued,
		IssuedAt:      &now,
		// Datos financieros copiados de la orden
		Subtotal:     order.Subtotal,
		Tax:          order.Tax,
		Discount:     order.Discount,
		ShippingCost: order.ShippingCost,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		// Datos del cliente
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

	// 4. Actualizar la orden con el invoice_id para que no aparezca como facturable
	if err := uc.repo.UpdateOrderInvoiceInfo(ctx, order.ID, fmt.Sprintf("%d", invoice.ID), ""); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("No se pudo actualizar invoice_id en la orden (la factura fue creada)")
	}

	uc.log.Info(ctx).
		Uint("invoice_id", invoice.ID).
		Str("invoice_number", invoice.InvoiceNumber).
		Str("order_id", order.ID).
		Msg("Factura manual registrada exitosamente")

	// 5. Publicar evento SSE para que el frontend se actualice en tiempo real
	if uc.ssePublisher != nil {
		invoice.OrderNumber = order.OrderNumber
		if err := uc.ssePublisher.PublishInvoiceCreated(ctx, invoice); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Error publicando evento SSE de factura manual")
		}
	}

	return invoice, nil
}
