package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

func (uc *useCase) DispatchOrderInvoicing(ctx context.Context, orderID string) (*entities.Invoice, error) {
	if uc.isInventoryExitOnly(ctx, orderID) {
		uc.log.Info(ctx).
			Str("order_id", orderID).
			Msg("Config en salida de inventario sin facturar: se emite comprobante contable en vez de factura")
		return uc.CreateJournal(ctx, &dtos.CreateJournalDTO{OrderID: orderID})
	}

	return uc.CreateInvoice(ctx, &dtos.CreateInvoiceDTO{OrderID: orderID, IsManual: false})
}

func (uc *useCase) isInventoryExitOnly(ctx context.Context, orderID string) bool {
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil || order == nil {
		return false
	}

	config, err := uc.repo.GetConfigByIntegration(ctx, order.IntegrationID)
	if err != nil {
		return false
	}
	if config == nil {
		config, err = uc.repo.GetEnabledConfigByBusiness(ctx, order.BusinessID)
		if err != nil {
			return false
		}
	}
	if config == nil || !config.Enabled || config.InvoiceConfig == nil {
		return false
	}

	inventoryExitOnly, _ := config.InvoiceConfig["inventory_exit_only"].(bool)
	if !inventoryExitOnly {
		return false
	}

	var integrationID uint
	if config.InvoicingIntegrationID != nil {
		integrationID = *config.InvoicingIntegrationID
	} else if config.InvoicingProviderID != nil {
		integrationID = *config.InvoicingProviderID
	} else {
		return false
	}

	provider, err := uc.resolveProvider(ctx, integrationID)
	if err != nil || provider != dtos.ProviderSiigo {
		uc.log.Warn(ctx).
			Str("order_id", orderID).
			Str("provider", provider).
			Msg("inventory_exit_only configurado en un proveedor que no es Siigo: se factura normal")
		return false
	}

	return true
}
